package handler

import (
	"context"
	"encoding/json"

	"github.com/xcheng85/session-monitor-k8s/internal/config"
	"github.com/xcheng85/session-monitor-k8s/internal/ddd"
	"github.com/xcheng85/session-monitor-k8s/internal/repository"
	"github.com/xcheng85/session-monitor-k8s/internal/session"
	"github.com/xcheng85/session-monitor-k8s/node/internal/domain"
	"go.uber.org/zap"
)

type domainEventHandlers[T ddd.IEvent] struct {
	logger         *zap.Logger
	config         config.IConfig
	repository     repository.IKVRepository // query server timestamp
	sessionService session.ISessionService  // set node ts cache
}

var _ ddd.IEventHandler[ddd.IEvent] = (*domainEventHandlers[ddd.IEvent])(nil)

func NewDomainEventHandlers(
	logger *zap.Logger,
	config config.IConfig,
	subscriber ddd.IEventDispatcher[ddd.IEvent],
	repository repository.IKVRepository,
	sessionService session.ISessionService,
) ddd.IEventHandler[ddd.IEvent] {
	handler := &domainEventHandlers[ddd.IEvent]{
		logger,
		config,
		repository,
		sessionService,
	}
	subscriber.Subscribe(handler,
		domain.NodeAddEvent,
		domain.NodeUpdateEvent,
		domain.NodeDeleteEvent,
		domain.NodeRecordNodeProvisionEvent,
		domain.NodeUpdateLabelsCacheEvent,
		domain.NodeInformerErrorEvent,
	)
	return handler
}

func (d domainEventHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
	d.logger.Sugar().Infof("HandleEvent: %s", event.EventName())
	switch event.EventName() {
	case domain.NodeAddEvent:
		return d.onNodeAdded(ctx, event)
	case domain.NodeUpdateEvent:
		return d.onNodeUpdated(ctx, event)
	case domain.NodeDeleteEvent:
		return d.onNodeDeleted(ctx, event)
	case domain.NodeUpdateLabelsCacheEvent:
		return d.onNodeUpdateLabelsCache(ctx, event)
	case domain.NodeRecordNodeProvisionEvent:
		return d.onRecordNodeProvisionTimestamp(ctx, event)
	}
	return nil
}

func (d domainEventHandlers[T]) onNodeAdded(ctx context.Context, event ddd.IEvent) error {
	payload := event.Payload().(*domain.NodeEventPayload)
	name, driverVersion := payload.Node.Name, payload.Node.DriverVersion
	d.logger.Sugar().Infow("Node is added", "Name", name, "DriverVersion", driverVersion)
	return nil
}

func (d domainEventHandlers[T]) onNodeDeleted(ctx context.Context, event ddd.IEvent) error {
	payload := event.Payload().(*domain.NodeEventPayload)
	name := payload.Node.Name
	d.logger.Sugar().Infow("Node is deleted", "Name", name)
	return nil
}

func (d domainEventHandlers[T]) onNodeUpdated(ctx context.Context, event ddd.IEvent) error {
	payload := event.Payload().(*domain.NodeEventPayload)
	name, driverVersion := payload.Node.Name, payload.Node.DriverVersion
	d.logger.Sugar().Infow("Node is updated", "Name", name, "DriverVersion", driverVersion)
	return nil
}

func (d domainEventHandlers[T]) onRecordNodeProvisionTimestamp(ctx context.Context, event ddd.IEvent) error {
	payload := event.Payload().(*domain.NodeEventPayload)
	nodeName := payload.Node.Name
	d.logger.Sugar().Infof("onRecordNodeProvisionTimestamp: %s", nodeName)
	// to do use repository to get server timestamp
	serverTimestamp, err := d.repository.GetServerTimestamp(ctx)
	if err != nil {
		d.logger.Sugar().Infof("repository.GetServerTimestamp has error: %s", err.Error())
	}

	d.logger.Sugar().Infof("Node %s is scheduled at %d", nodeName, serverTimestamp)
	err = d.sessionService.SetNodeProvisionTimeStamp(&session.SetNodeProvisionTimeStampActionPayload{
		NodeName:  nodeName,
		Timestamp: serverTimestamp,
	})
	return err
}

func (d domainEventHandlers[T]) onNodeUpdateLabelsCache(ctx context.Context, event ddd.IEvent) error {
	payload := event.Payload().(*domain.NodeEventPayload)
	nodeLables := payload.Node.Labels
	agentPoolName := (*nodeLables)["agentpool"]
	j, err := json.Marshal(nodeLables)
	if err != nil {
		d.logger.Sugar().Errorf("agentPoolName %s json marshal node labels has error: %s", agentPoolName, err.Error())
		return domain.NewBadNodeLabelErr(nodeLables)
	}
	gpuAgentPoolSetKey := d.config.Get("app.gpu_agent_pool_set_key").(string)
	// transaction
	var numKeysAdded int64
	numKeysAdded, err = d.repository.AddToUnsortedSet(ctx, gpuAgentPoolSetKey, &repository.Object{
		Key:     agentPoolName,
		Payload: string(j),
	})
	if err != nil {
		d.logger.Sugar().Errorf("AddToUnsortedSet has error: %s", err.Error())
	}
	d.logger.Sugar().Infof("AddToUnsortedSet: %d key(s) are added", numKeysAdded)
	return err
}
