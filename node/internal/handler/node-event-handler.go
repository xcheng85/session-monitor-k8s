package handler

import (
	"context"
	"errors"

	"github.com/xcheng85/session-monitor-k8s/internal/config"
	"github.com/xcheng85/session-monitor-k8s/internal/ddd"
	"github.com/xcheng85/session-monitor-k8s/internal/k8s"
	"github.com/xcheng85/session-monitor-k8s/node/internal/domain"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
)

type NodeLabel string

const (
	NVIDIA_DRIVER_VERSION_MAJOR NodeLabel = "nvidia.com/cuda.driver.major"
	NVIDIA_DRIVER_VERSION_MINOR NodeLabel = "nvidia.com/cuda.driver.minor"
	NVIDIA_DRIVER_VERSION_REV   NodeLabel = "nvidia.com/cuda.driver.rev"
)

type NodeEventHandler struct {
	ctx                   context.Context
	logger                *zap.Logger
	config                config.IConfig
	domainEventDispatcher ddd.IEventDispatcher[ddd.IEvent]
}

func NewNodeEventHandler(
	ctx context.Context,
	logger *zap.Logger,
	config config.IConfig,
	domainEventDispatcher ddd.IEventDispatcher[ddd.IEvent],
	domainEventHandler ddd.IEventHandler[ddd.IEvent]) k8s.IK8sEventHandler {
	return &NodeEventHandler{
		ctx,
		logger,
		config,
		domainEventDispatcher,
	}
}

func (handler *NodeEventHandler) CustomWatchErrorHandler(r *cache.Reflector, err error) {
	handler.logger.Sugar().Errorw("Watch error", err)
	handler.domainEventDispatcher.Publish(handler.ctx, ddd.NewEvent(
		domain.NodeInformerErrorEvent,
		&domain.NodeInformerErrorPayload{
			Err: err,
		},
	))
}

func (handler *NodeEventHandler) shouldIgnore(nodeLables *map[string]string, gpuObserveeMap *map[string]string) bool {
	nodeShouldBeIgnored := false
	for key, value := range *gpuObserveeMap {
		value2, exist := (*nodeLables)[key]
		if !exist || value2 != value {
			nodeShouldBeIgnored = true
			break
		}
	}
	return nodeShouldBeIgnored
}

func (handler *NodeEventHandler) OnAddObject(obj interface{}) {
	node, err := parseNode(obj.(*unstructured.Unstructured))
	if err != nil {
		handler.logger.Sugar().Error("OnAddObject:", err)
	} else {
		name := node.Name
		agentPoolName := node.Labels["agentpool"]
		gpu_observee_labels_interface := handler.config.Get("app.gpu_observee_labels").([]interface{})
		gpuObserveeMap := map[string]string{}
		for i := 0; i < len(gpu_observee_labels_interface)/2; i++ {
			k, v := gpu_observee_labels_interface[2*i].(string), gpu_observee_labels_interface[2*i+1].(string)
			gpuObserveeMap[k] = v
		}
		if !handler.shouldIgnore(&node.Labels, &gpuObserveeMap) {
			handler.logger.Sugar().Infof("[OnAddObject] observer agentPoolName: %f", agentPoolName)
			driverVersion, _ := parseGPUDriverVersion(&node.Labels, handler.logger)
			// two event to dispatch
			// 1. updateGPUNodeAgentPoolLabelsCache
			// 2. record node provision timestamp
			nodeDomain := &domain.Node{
				Name:          name,
				DriverVersion: driverVersion,
				Labels:        &node.Labels,
			}
			handler.domainEventDispatcher.Publish(
				handler.ctx,
				ddd.NewEvent(
					domain.NodeAddEvent,
					&domain.NodeEventPayload{
						Node: nodeDomain,
					}),
				ddd.NewEvent(
					domain.NodeUpdateLabelsCacheEvent,
					&domain.NodeEventPayload{
						Node: nodeDomain,
					}),
				ddd.NewEvent(
					domain.NodeRecordNodeProvisionEvent,
					&domain.NodeEventPayload{
						Node: nodeDomain,
					},
				),
			)
		}
	}
}

func (handler *NodeEventHandler) OnUpdateObject(oldObj, newObj interface{}) {
	node, err := parseNode(newObj.(*unstructured.Unstructured))
	if err != nil {
		handler.logger.Sugar().Error("OnUpdateObject:", err)
	} else {
		name := node.Name
		agentPoolName := node.Labels["agentpool"]
		gpu_observee_labels_interface := handler.config.Get("app.gpu_observee_labels").([]interface{})
		gpuObserveeMap := map[string]string{}
		for i := 0; i < len(gpu_observee_labels_interface)/2; i++ {
			k, v := gpu_observee_labels_interface[2*i].(string), gpu_observee_labels_interface[2*i+1].(string)
			gpuObserveeMap[k] = v
		}
		if !handler.shouldIgnore(&node.Labels, &gpuObserveeMap) {
			handler.logger.Sugar().Infof("[OnDeleteObject] observer agentPoolName: %f", agentPoolName)
			driverVersion, driverVersionErr := parseGPUDriverVersion(&node.Labels, handler.logger)
			// one event to dispatch
			// 1. updateGPUNodeAgentPoolLabelsCache only driverVersion is ready
			nodeDomain := &domain.Node{
				Name:          name,
				DriverVersion: driverVersion,
				Labels:        &node.Labels,
			}
			if driverVersionErr == nil {
				handler.domainEventDispatcher.Publish(
					handler.ctx,
					ddd.NewEvent(
						domain.NodeUpdateEvent,
						&domain.NodeEventPayload{
							Node: nodeDomain,
						}),
					ddd.NewEvent(
						domain.NodeUpdateLabelsCacheEvent,
						&domain.NodeEventPayload{
							Node: nodeDomain,
						}))
			} else {
				handler.domainEventDispatcher.Publish(
					handler.ctx,
					ddd.NewEvent(
						domain.NodeUpdateEvent,
						&domain.NodeEventPayload{
							Node: nodeDomain,
						}))
			}
		}
	}
}

func (handler *NodeEventHandler) OnDeleteObject(obj interface{}) {
	node, err := parseNode(obj.(*unstructured.Unstructured))
	if err != nil {
		handler.logger.Sugar().Error("OnDeleteObject:", err)
	} else {
		name := node.Name
		agentPoolName := node.Labels["agentpool"]
		gpu_observee_labels_interface := handler.config.Get("app.gpu_observee_labels").([]interface{})
		gpuObserveeMap := map[string]string{}
		for i := 0; i < len(gpu_observee_labels_interface)/2; i++ {
			k, v := gpu_observee_labels_interface[2*i].(string), gpu_observee_labels_interface[2*i+1].(string)
			gpuObserveeMap[k] = v
		}
		handler.logger.Sugar().Info(gpuObserveeMap)
		if !handler.shouldIgnore(&node.Labels, &gpuObserveeMap) {
			handler.logger.Sugar().Infof("[OnDeleteObject] observer agentPoolName: %f", agentPoolName)
			nodeDomain := &domain.Node{
				Name:   name,
				Labels: &node.Labels,
			}
			handler.domainEventDispatcher.Publish(
				handler.ctx,
				ddd.NewEvent(
					domain.NodeDeleteEvent,
					&domain.NodeEventPayload{
						Node: nodeDomain,
					}))

		}
	}
}

func parseNode(u *unstructured.Unstructured) (node v1.Node, err error) {
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &node)
	return node, err
}

func parseGPUDriverVersion(nodeLables *map[string]string, logger *zap.Logger) (string, error) {
	driverVersionMajor, driverVersionMajorExist := (*nodeLables)[string(NVIDIA_DRIVER_VERSION_MAJOR)]
	driverVersionMinor, driverVersionMinorExist := (*nodeLables)[string(NVIDIA_DRIVER_VERSION_MINOR)]
	driverVersionRev, driverVersionRevExist := (*nodeLables)[string(NVIDIA_DRIVER_VERSION_REV)]
	logger.Sugar().Info(nodeLables)
	logger.Sugar().Info(driverVersionMajor)
	logger.Sugar().Info(driverVersionMinor)
	logger.Sugar().Info(driverVersionRev)
	if driverVersionMajorExist && driverVersionMinorExist && driverVersionRevExist {
		driverVersion := driverVersionMajor + "." + driverVersionMinor + "." + driverVersionRev
		return driverVersion, nil
	} else {
		return "", errors.New("Driver Version is missing!")
	}
}
