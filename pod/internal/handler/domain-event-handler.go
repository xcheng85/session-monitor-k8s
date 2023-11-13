package handler

import (
	"context"
	"github.com/xcheng85/session-monitor-k8s/internal/ddd"
	"github.com/xcheng85/session-monitor-k8s/internal/repository"
	"github.com/xcheng85/session-monitor-k8s/internal/session"
	"github.com/xcheng85/session-monitor-k8s/pod/internal/domain"
	"go.uber.org/zap"
)

type domainEventHandlers[T ddd.IEvent] struct {
	logger         *zap.Logger
	repository     repository.IKVRepository
	sessionService session.ISessionService
}

var _ ddd.IEventHandler[ddd.IEvent] = (*domainEventHandlers[ddd.IEvent])(nil)

func NewDomainEventHandlers(logger *zap.Logger,
	subscriber ddd.IEventDispatcher[ddd.IEvent],
	repository repository.IKVRepository,
	sessionService session.ISessionService,
) ddd.IEventHandler[ddd.IEvent] {
	handler := &domainEventHandlers[ddd.IEvent]{
		logger,
		repository,
		sessionService,
	}
	subscriber.Subscribe(handler,
		domain.PodAddEvent,
		domain.PodDeleteEvent,
		domain.PodReadyEvent,
		domain.PodRecordPodScheduleEvent,
	)
	return handler
}

func (d domainEventHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
	switch event.EventName() {
	case domain.PodAddEvent:
		return d.onPodAdded(ctx, event)
	case domain.PodDeleteEvent:
		return d.onPodDeleted(ctx, event)
	case domain.PodRecordPodScheduleEvent:
		return d.onRecordPodScheduleTimestamp(ctx, event)
	case domain.PodReadyEvent:
		return d.onPodReady(ctx, event)
	}
	return nil
}

func (d domainEventHandlers[T]) onPodAdded(ctx context.Context, event ddd.IEvent) error {
	payload := event.Payload().(*domain.PodEventPayload)
	name, namespace, sessionId := payload.Pod.Name, payload.Pod.Namespace, payload.Pod.SessionId
	d.logger.Sugar().Infow("Pod is added", "Name", name, "Namespace", namespace, "SessionId", sessionId)
	return nil
}

func (d domainEventHandlers[T]) onPodDeleted(ctx context.Context, event ddd.IEvent) error {
	payload := event.Payload().(*domain.PodEventPayload)
	name, namespace, sessionId := payload.Pod.Name, payload.Pod.Namespace, payload.Pod.SessionId
	d.logger.Sugar().Infow("Pod is deleted", "Name", name, "Namespace", namespace, "SessionId", sessionId)
	err := d.sessionService.SetSessionDeletable(&session.SetSessionDeletableActionPayload{
		SessionId: sessionId,
	})
	return err
}

func (d domainEventHandlers[T]) onRecordPodScheduleTimestamp(ctx context.Context, event ddd.IEvent) error {
	payload := event.Payload().(*domain.PodEventPayload)
	sessionId := payload.Pod.SessionId
	// to do use repository to get server timestamp
	d.logger.Sugar().Infof("Session %s is scheduled at %d", sessionId, 0)
	serverTimestamp, err := d.repository.GetServerTimestamp(ctx)
	if err == nil {
		err = d.sessionService.SetPodScheduleTimeStamp(&session.UpdateSessionTimeStampLikeFieldActionPayload{
			SessionId: sessionId,
			Timestamp: serverTimestamp,
		})
	}
	return err
}

func (d domainEventHandlers[T]) onPodReady(ctx context.Context, event ddd.IEvent) error {
	payload := event.Payload().(*domain.PodEventPayload)
	name, namespace, sessionId, nodeName, ip := payload.Pod.Name, payload.Pod.Namespace, payload.Pod.SessionId,
		payload.Pod.NodeName, payload.Pod.Ip
	d.logger.Sugar().Infow("Pod should be ready and enqueued", "Name", name,
		"Namespace", namespace,
		"SessionId", sessionId,
		"NodeName", nodeName,
		"Ip", ip,
	)

	nodeProvisionedTimeStamp, nodeErr := d.sessionService.GetNodeProvisionTimeStamp(nodeName)
	podScheduledTimeStamp, podErr := d.sessionService.GetPodScheduleTimeStamp(sessionId)

	if nodeErr != nil {
		return nodeErr
	}

	if podErr != nil {
		return podErr
	}

	return d.sessionService.SetSessionReady(&session.SetSessionReadyActionPayload{
		SessionId:              sessionId,
		NodeName:               nodeName,
		PodInternalIp:          ip,
		NodeProvisionTimeStamp: nodeProvisionedTimeStamp,
		PodScheduleTimeStamp:   podScheduledTimeStamp,
	})
}
