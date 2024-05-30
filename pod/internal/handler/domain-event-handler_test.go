package handler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	// "github.com/stretchr/testify/mock"
	// "github.com/xcheng85/session-monitor-k8s/internal/config"
	"github.com/xcheng85/session-monitor-k8s/internal/ddd"
	"github.com/xcheng85/session-monitor-k8s/internal/logger"
	"github.com/xcheng85/session-monitor-k8s/internal/repository"
	"github.com/xcheng85/session-monitor-k8s/internal/session"
	// "github.com/xcheng85/session-monitor-k8s/internal/repository"
	// "github.com/xcheng85/session-monitor-k8s/internal/session"
	"github.com/xcheng85/session-monitor-k8s/pod/internal/domain"
	// "go.uber.org/zap"
)

var (
	nodeName      = "NodeName"
	sessionId     = "SessionId"
	podInternalIp = "8.8.8.8"
	pod           = &domain.Pod{
		Name:      "Name",
		Namespace: "Namespace",
		SessionId: sessionId,
		NodeName:  nodeName,
		Ip:        podInternalIp,
	}

	serverTimestamp          = int64(8888888888888)
	nodeProvisionedTimestamp = int64(6888888888888)
	podScheduledTimestamp    = int64(7888888888888)
)

func TestNewDomainEventHandlers(t *testing.T) {
	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})
	mockEventDispatcher := &ddd.MockIEventDispatcher[ddd.IEvent]{}
	mockEventDispatcher.On("Subscribe", mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockKVRepository := &repository.MockIKVRepository{}
	mockSessionService := &session.MockISessionService{}

	h := NewDomainEventHandlers(logger, mockEventDispatcher, mockKVRepository, mockSessionService)
	assert.NotNil(t, h, "Pod Event Handler should not be nil")
}

func TestHandleEvent_PodAddEvent(t *testing.T) {
	ctx := context.TODO()
	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})
	mockEventDispatcher := &ddd.MockIEventDispatcher[ddd.IEvent]{}
	mockEventDispatcher.On("Subscribe", mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockKVRepository := &repository.MockIKVRepository{}
	mockSessionService := &session.MockISessionService{}

	h := NewDomainEventHandlers(logger, mockEventDispatcher, mockKVRepository, mockSessionService)
	err := h.HandleEvent(ctx, ddd.NewEvent(
		domain.PodAddEvent,
		&domain.PodEventPayload{
			Pod: pod,
		}))
	assert.Nil(t, err, "Handle PodAddEvent should not throw err")
}

func TestHandleEvent_PodDeleted(t *testing.T) {
	ctx := context.TODO()
	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})
	mockEventDispatcher := &ddd.MockIEventDispatcher[ddd.IEvent]{}
	mockEventDispatcher.On("Subscribe", mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockKVRepository := &repository.MockIKVRepository{}
	mockSessionService := &session.MockISessionService{}
	mockSessionService.On("SetSessionDeletable", &session.SetSessionDeletableActionPayload{
		SessionId: "SessionId",
		CallerId:  "Session-monitor-service",
	}).Return(nil).Once()

	h := NewDomainEventHandlers(logger, mockEventDispatcher, mockKVRepository, mockSessionService)
	err := h.HandleEvent(ctx, ddd.NewEvent(
		domain.PodDeleteEvent,
		&domain.PodEventPayload{
			Pod: pod,
		}))
	assert.Nil(t, err, "Handle PodDeleted Event should not throw err")
}

func TestHandleEvent_RecordPodScheduleTimestamp(t *testing.T) {
	ctx := context.TODO()
	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})
	mockEventDispatcher := &ddd.MockIEventDispatcher[ddd.IEvent]{}
	mockEventDispatcher.On("Subscribe", mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockKVRepository := &repository.MockIKVRepository{}
	mockKVRepository.On("GetServerTimestamp", ctx).Return(serverTimestamp, nil)
	mockSessionService := &session.MockISessionService{}
	mockSessionService.On("SetPodScheduleTimeStamp", &session.UpdateSessionTimeStampLikeFieldActionPayload{
		SessionId: "SessionId",
		Timestamp: serverTimestamp,
	}).Return(nil).Once()
	h := NewDomainEventHandlers(logger, mockEventDispatcher, mockKVRepository, mockSessionService)
	err := h.HandleEvent(ctx, ddd.NewEvent(
		domain.PodRecordPodScheduleEvent,
		&domain.PodEventPayload{
			Pod: pod,
		}))
	assert.Nil(t, err, "Handle RecordPodScheduleTimestamp Event should not throw err")
}

func TestHandleEvent_PodReady(t *testing.T) {
	ctx := context.TODO()
	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})
	mockEventDispatcher := &ddd.MockIEventDispatcher[ddd.IEvent]{}
	mockEventDispatcher.On("Subscribe", mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockKVRepository := &repository.MockIKVRepository{}
	mockSessionService := &session.MockISessionService{}
	mockSessionService.On("GetNodeProvisionTimeStamp", nodeName).Return(nodeProvisionedTimestamp, nil)
	mockSessionService.On("GetPodScheduleTimeStamp", sessionId).Return(podScheduledTimestamp, nil)
	mockSessionService.On("SetSessionReady", &session.SetSessionReadyActionPayload{
		SessionId:              sessionId,
		NodeName:               nodeName,
		PodInternalIp:          podInternalIp,
		NodeProvisionTimeStamp: nodeProvisionedTimestamp,
		PodScheduleTimeStamp:   podScheduledTimestamp,
	}).Return(nil)
	h := NewDomainEventHandlers(logger, mockEventDispatcher, mockKVRepository, mockSessionService)
	err := h.HandleEvent(ctx, ddd.NewEvent(
		domain.PodReadyEvent,
		&domain.PodEventPayload{
			Pod: pod,
		}))
	assert.Nil(t, err, "Handle PodReady Event should not throw err")
	mockSessionService.AssertNumberOfCalls(t, "GetNodeProvisionTimeStamp", 1)
	mockSessionService.AssertNumberOfCalls(t, "GetPodScheduleTimeStamp", 1)
	mockSessionService.AssertNumberOfCalls(t, "SetSessionReady", 1)
}

func TestHandleEvent_PodReady_Negative_GetNodeProvisionTimeStamp(t *testing.T) {
	ctx := context.TODO()
	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})
	mockEventDispatcher := &ddd.MockIEventDispatcher[ddd.IEvent]{}
	mockEventDispatcher.On("Subscribe", mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockKVRepository := &repository.MockIKVRepository{}
	mockSessionService := &session.MockISessionService{}
	mockSessionService.On("GetNodeProvisionTimeStamp", nodeName).Return(int64(0), session.NewInvalidStoreKeyErr("bogus"))
	mockSessionService.On("GetPodScheduleTimeStamp", sessionId).Return(podScheduledTimestamp, nil)
	mockSessionService.On("SetSessionReady", &session.SetSessionReadyActionPayload{
		SessionId:              sessionId,
		NodeName:               nodeName,
		PodInternalIp:          podInternalIp,
		NodeProvisionTimeStamp: nodeProvisionedTimestamp,
		PodScheduleTimeStamp:   podScheduledTimestamp,
	}).Return(nil)
	h := NewDomainEventHandlers(logger, mockEventDispatcher, mockKVRepository, mockSessionService)
	err := h.HandleEvent(ctx, ddd.NewEvent(
		domain.PodReadyEvent,
		&domain.PodEventPayload{
			Pod: pod,
		}))
	assert.Equal(t, session.NewInvalidStoreKeyErr("bogus"), err, "Handle PodReady Event should not throw err")
	mockSessionService.AssertNumberOfCalls(t, "GetNodeProvisionTimeStamp", 1)
	mockSessionService.AssertNumberOfCalls(t, "GetPodScheduleTimeStamp", 0)
	mockSessionService.AssertNumberOfCalls(t, "SetSessionReady", 0)
}

func TestHandleEvent_PodReady_Negative_GetPodScheduleTimeStamp(t *testing.T) {
	ctx := context.TODO()
	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})
	mockEventDispatcher := &ddd.MockIEventDispatcher[ddd.IEvent]{}
	mockEventDispatcher.On("Subscribe", mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockKVRepository := &repository.MockIKVRepository{}
	mockSessionService := &session.MockISessionService{}
	mockSessionService.On("GetNodeProvisionTimeStamp", nodeName).Return(nodeProvisionedTimestamp, nil)
	mockSessionService.On("GetPodScheduleTimeStamp", sessionId).Return(int64(0), session.NewInvalidStoreKeyErr("bogus"))
	mockSessionService.On("SetSessionReady", &session.SetSessionReadyActionPayload{
		SessionId:              sessionId,
		NodeName:               nodeName,
		PodInternalIp:          podInternalIp,
		NodeProvisionTimeStamp: nodeProvisionedTimestamp,
		PodScheduleTimeStamp:   podScheduledTimestamp,
	}).Return(nil)
	h := NewDomainEventHandlers(logger, mockEventDispatcher, mockKVRepository, mockSessionService)
	err := h.HandleEvent(ctx, ddd.NewEvent(
		domain.PodReadyEvent,
		&domain.PodEventPayload{
			Pod: pod,
		}))
	assert.Equal(t, session.NewInvalidStoreKeyErr("bogus"), err, "Handle PodReady Event should not throw err")
	mockSessionService.AssertNumberOfCalls(t, "GetNodeProvisionTimeStamp", 1)
	mockSessionService.AssertNumberOfCalls(t, "GetPodScheduleTimeStamp", 1)
	mockSessionService.AssertNumberOfCalls(t, "SetSessionReady", 0)
}
