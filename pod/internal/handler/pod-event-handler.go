package handler

import (
	"context"
	"encoding/json"

	"github.com/thoas/go-funk"
	"github.com/xcheng85/session-monitor-k8s/internal/ddd"
	"github.com/xcheng85/session-monitor-k8s/internal/k8s"
	"github.com/xcheng85/session-monitor-k8s/pod/internal/domain"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
)

type PodEventHandler struct {
	ctx                   context.Context
	logger                *zap.Logger
	domainEventDispatcher ddd.IEventDispatcher[ddd.IEvent]
}

func NewPodEventHandler(ctx context.Context, logger *zap.Logger,
	domainEventDispatcher ddd.IEventDispatcher[ddd.IEvent],
	domainEventHandler ddd.IEventHandler[ddd.IEvent],
) k8s.IK8sEventHandler {
	return &PodEventHandler{
		ctx,
		logger,
		domainEventDispatcher,
	}
}

func (handler *PodEventHandler) CustomWatchErrorHandler(r *cache.Reflector, err error) {
	handler.logger.Sugar().Errorw("Watch error", err)
	handler.domainEventDispatcher.Publish(handler.ctx, ddd.NewEvent(
		domain.PodInformerErrorEvent,
		&domain.PodInformerErrorPayload{
			Err: err,
		},
	))
}

func (handler *PodEventHandler) OnAddObject(obj interface{}) {
	pod, err := parsePod(obj.(*unstructured.Unstructured))
	if err != nil {
		handler.logger.Sugar().Error("OnAddObject:", err)
	} else {
		name, namespace, sessionId := pod.ObjectMeta.Name, pod.ObjectMeta.Namespace, pod.ObjectMeta.Labels["sessionId"]
		handler.domainEventDispatcher.Publish(handler.ctx, ddd.NewEvent(
			domain.PodAddEvent,
			&domain.PodEventPayload{
				Pod: &domain.Pod{
					Name:      name,
					Namespace: namespace,
					SessionId: sessionId,
				},
			},
		))
	}
}

func (handler *PodEventHandler) OnUpdateObject(oldObj, newObj interface{}) {
	eventName := domain.PodNilEvent
	var eventPlayload interface{}

	pod, err := parsePod(newObj.(*unstructured.Unstructured))
	if err != nil {
		handler.logger.Sugar().Error("OnUpdateObject:", err)
	}
	// label managed allows dev's smoke test, which is living outside of session management backend
	name, namespace, sessionId, isManaged, phase, nodeName, conditions, ip :=
		pod.ObjectMeta.Name, pod.ObjectMeta.Namespace, pod.ObjectMeta.Labels["sessionId"],
		pod.ObjectMeta.Labels["managed"], pod.Status.Phase, pod.Spec.NodeName,
		pod.Status.Conditions, pod.Status.PodIP

	m := funk.ToMap(conditions, "Type").(map[v1.PodConditionType]v1.PodCondition)
	if sessionId != "" && isManaged != "false" {
		handler.logger.Sugar().Infow("Pod is updated", "Name", name, "Namespace", namespace, "SessionId", sessionId, "Phase", pod.Status.Phase, "PodIP:", pod.Status.PodIP)
		if phase == v1.PodFailed || phase == v1.PodSucceeded || phase == v1.PodUnknown {
			handler.logger.Sugar().Infow("Pod should be deleted", "Name", name, "Namespace", namespace, "SessionId", sessionId)
			eventName = domain.PodDeleteEvent
			eventPlayload = &domain.PodEventPayload{
				Pod: &domain.Pod{
					Name:      name,
					Namespace: namespace,
					SessionId: sessionId,
				},
			}
		} else if phase == v1.PodPending && m[v1.PodScheduled].Status == v1.ConditionTrue {
			eventName = domain.PodRecordPodScheduleEvent
			eventPlayload = &domain.PodEventPayload{
				Pod: &domain.Pod{
					Name:      name,
					Namespace: namespace,
					SessionId: sessionId,
				},
			}
		} else if phase == v1.PodRunning {
			conditions := pod.Status.Conditions
			m := funk.ToMap(conditions, "Type").(map[v1.PodConditionType]v1.PodCondition)
			handler.logger.Sugar().Infow("Running Pod Status Update", "Name", name, "Namespace",
				namespace, "SessionId", sessionId, "PodInitialized",
				m[v1.PodInitialized].Status, "PodScheduled", m[v1.PodScheduled].Status, "ContainersReady",
				m[v1.ContainersReady].Status, "PodReady", m[v1.PodReady].Status)

			if m[v1.PodInitialized].Status == v1.ConditionTrue && m[v1.PodScheduled].Status == v1.ConditionTrue && m[v1.ContainersReady].Status == v1.ConditionTrue && m[v1.PodReady].Status == v1.ConditionTrue {
				if pod.ObjectMeta.DeletionTimestamp == nil {
					eventName = domain.PodReadyEvent
					eventPlayload = &domain.PodEventPayload{
						Pod: &domain.Pod{
							Name:      name,
							Namespace: namespace,
							SessionId: sessionId,
							NodeName:  nodeName,
							Ip:        ip,
						},
					}
				} else {
					handler.logger.Sugar().Infow("Pod should be deleted", "Name", name, "Namespace", namespace, "SessionId", sessionId)
					eventName = domain.PodDeleteEvent
					eventPlayload = &domain.PodEventPayload{
						Pod: &domain.Pod{
							Name:      name,
							Namespace: namespace,
							SessionId: sessionId,
						},
					}
				}
			}
			// logic to catch crashed containers
			r := funk.Filter(pod.Status.ContainerStatuses, func(c v1.ContainerStatus) bool {
				return c.State.Terminated != nil
			}).([]v1.ContainerStatus)
			if len(r) > 0 {
				handler.logger.Sugar().Infow("Pod has crashed containers", "Name", name, "Namespace", namespace, "SessionId", sessionId)
				handler.logger.Sugar().Infow("Pod should be deleted", "Name", name, "Namespace", namespace, "SessionId", sessionId)
				out, _ := json.Marshal(r)
				handler.logger.Sugar().Infow(string(out))

				eventName = domain.PodDeleteEvent
				eventPlayload = &domain.PodEventPayload{
					Pod: &domain.Pod{
						Name:      name,
						Namespace: namespace,
						SessionId: sessionId,
					},
				}
			}
		}
	}
	if eventName != domain.PodNilEvent {
		handler.domainEventDispatcher.Publish(handler.ctx, ddd.NewEvent(
			eventName,
			eventPlayload,
		))
	}
}

func (handler *PodEventHandler) OnDeleteObject(obj interface{}) {
	pod, err := parsePod(obj.(*unstructured.Unstructured))
	if err != nil {
		handler.logger.Sugar().Error("OnAddObject:", err)
	} else {
		name, namespace := pod.ObjectMeta.Name, pod.ObjectMeta.Namespace
		handler.logger.Sugar().Infow("Pod is deleted", "Name", name, "Namespace", namespace)
	}
}

func parsePod(u *unstructured.Unstructured) (pod v1.Pod, err error) {
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &pod)
	return pod, err
}
