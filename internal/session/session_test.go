package session

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/xcheng85/session-monitor-k8s/internal/config"
	"github.com/xcheng85/session-monitor-k8s/internal/logger"
	"github.com/xcheng85/session-monitor-k8s/internal/repository"
)

func TestSetSessionReady(t *testing.T) {
	mockEnqueueSessionStreamKey := "enqueue_session_stream_key"
	mockKVRepository := &repository.MockIKVRepository{}
	mockServerTimestamp, mockNodeProvisionTimeStamp, mockPodScheduleTimestamp := int64(88888888888),
		int64(188888888888), int64(28888888888)

	ctx := context.TODO()
	mockConfig := &config.MockIConfig{}
	mockConfig.On("Get", "app.enqueue_session_stream_key").Return(mockEnqueueSessionStreamKey, nil).Once()
	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})
	mockKVRepository.On("GetServerTimestamp", ctx).Return(mockServerTimestamp, nil).Once()
	mockKVRepository.On("AddStreamEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("streamId-1", nil).Once()
	mockPayload := SetSessionReadyActionPayload{
		SessionId:              "sessionId",
		NodeName:               "nodeName",
		PodInternalIp:          "8.8.8.8",
		NodeProvisionTimeStamp: mockNodeProvisionTimeStamp,
		PodScheduleTimeStamp:   mockPodScheduleTimestamp,
	}
	mockPayloadBuf, _ := json.Marshal(mockPayload)

	sessionService := NewSessionService(ctx, logger, mockConfig, mockKVRepository)
	err := sessionService.SetSessionReady(&mockPayload)
	assert.Nil(t, err, "sessionService.SetSessionReady should not throw error")
	mockKVRepository.AssertNumberOfCalls(t, "GetServerTimestamp", 1)
	mockKVRepository.AssertNumberOfCalls(t, "AddStreamEvent", 1)
	mockKVRepository.AssertCalled(t, "AddStreamEvent", ctx, mockEnqueueSessionStreamKey, "*", []interface{}{"TaskType",
		string(EnqueueSession), "TaskInfo", string(mockPayloadBuf),
		"TaskCreateTimeStamp", mockServerTimestamp})
}

func TestSetSessionDeletable(t *testing.T) {
	mockDeleteSessionStreamKey := "delete_session_stream_key"
	mockKVRepository := &repository.MockIKVRepository{}
	mockServerTimestamp := int64(88888888888)

	ctx := context.TODO()
	mockConfig := &config.MockIConfig{}
	mockConfig.On("Get", "app.delete_session_stream_key").Return(mockDeleteSessionStreamKey, nil).Once()
	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})
	mockKVRepository.On("GetServerTimestamp", ctx).Return(mockServerTimestamp, nil).Once()
	mockKVRepository.On("AddStreamEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("streamId-1", nil).Once()
	mockPayload := SetSessionDeletableActionPayload{
		SessionId: "sessionId",
	}
	mockPayloadBuf, _ := json.Marshal(mockPayload)

	sessionService := NewSessionService(ctx, logger, mockConfig, mockKVRepository)
	err := sessionService.SetSessionDeletable(&mockPayload)
	assert.Nil(t, err, "sessionService.SetSessionDeletable should not throw error")
	mockKVRepository.AssertNumberOfCalls(t, "GetServerTimestamp", 1)
	mockKVRepository.AssertNumberOfCalls(t, "AddStreamEvent", 1)
	mockKVRepository.AssertCalled(t, "AddStreamEvent", ctx, mockDeleteSessionStreamKey, "*", []interface{}{"TaskType",
		string(DeleteSession), "TaskInfo", string(mockPayloadBuf),
		"TaskCreateTimeStamp", mockServerTimestamp})
}

func TestSetNodeProvisionTimeStamp(t *testing.T) {
	mockKVRepository := &repository.MockIKVRepository{}
	mockNodeName := "nodeName"
	mockNodeProvisioningTimestamp := int64(88888888888)

	ctx := context.TODO()
	mockConfig := &config.MockIConfig{}
	mockConfig.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})

	mockPayload := SetNodeProvisionTimeStampActionPayload{
		NodeName:  mockNodeName,
		Timestamp: mockNodeProvisioningTimestamp,
	}
	sessionService := NewSessionService(ctx, logger, mockConfig, mockKVRepository)
	err := sessionService.SetNodeProvisionTimeStamp(&mockPayload)
	assert.Nil(t, err, "sessionService.SetNodeProvisionTimeStamp should not throw error")

	mockConfig.AssertNumberOfCalls(t, "Set", 1)
	mockConfig.AssertCalled(t, "Set", "NodeProvisionTimeStamp.nodeName", mockNodeProvisioningTimestamp)
}

func TestSetPodScheduleTimeStamp(t *testing.T) {
	mockKVRepository := &repository.MockIKVRepository{}
	mockSessionId := "sessionId"
	mockSetPodScheduleTimeStamp := int64(88888888888)

	ctx := context.TODO()
	mockConfig := &config.MockIConfig{}
	mockConfig.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})

	mockPayload := UpdateSessionTimeStampLikeFieldActionPayload{
		SessionId: mockSessionId,
		Timestamp: mockSetPodScheduleTimeStamp,
	}
	sessionService := NewSessionService(ctx, logger, mockConfig, mockKVRepository)
	err := sessionService.SetPodScheduleTimeStamp(&mockPayload)
	assert.Nil(t, err, "sessionService.SetPodScheduleTimeStamp should not throw error")

	mockConfig.AssertNumberOfCalls(t, "Set", 1)
	mockConfig.AssertCalled(t, "Set", "PodScheduleTimeStamp.sessionId", mockSetPodScheduleTimeStamp)
}

func TestGetNodeProvisionTimeStamp(t *testing.T) {
	mockKVRepository := &repository.MockIKVRepository{}
	mockNodeName := "nodeName"
	mockNodeProvisioningTimestamp := int64(88888888888)

	ctx := context.TODO()
	mockConfig := &config.MockIConfig{}
	mockConfig.On("Get", mock.Anything).Return(mockNodeProvisioningTimestamp).Once()
	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})
	sessionService := NewSessionService(ctx, logger, mockConfig, mockKVRepository)
	timestamp, err := sessionService.GetNodeProvisionTimeStamp(mockNodeName)
	assert.Nil(t, err, "sessionService.GetNodeProvisionTimeStamp should not throw error")
	assert.Equal(t, mockNodeProvisioningTimestamp, timestamp)
	mockConfig.AssertNumberOfCalls(t, "Get", 1)
	mockConfig.AssertCalled(t, "Get", "NodeProvisionTimeStamp.nodeName")
}

func TestGetNodeProvisionTimeStampKeyNotFound(t *testing.T) {
	mockKVRepository := &repository.MockIKVRepository{}
	mockNodeName := "nodeName"
	ctx := context.TODO()
	mockConfig := &config.MockIConfig{}
	mockConfig.On("Get", mock.Anything).Return(nil).Once()
	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})
	sessionService := NewSessionService(ctx, logger, mockConfig, mockKVRepository)
	_, err := sessionService.GetNodeProvisionTimeStamp(mockNodeName)
	assert.True(t, NewInvalidStoreKeyErr("NodeProvisionTimeStamp.nodeName").Is(err))
}

func TestGetPodScheduleTimeStamp(t *testing.T) {
	mockKVRepository := &repository.MockIKVRepository{}
	mockSessionId := "sessionId"
	mockPodScheduleTimestamp := int64(88888888888)

	ctx := context.TODO()
	mockConfig := &config.MockIConfig{}
	mockConfig.On("Get", mock.Anything).Return(mockPodScheduleTimestamp).Once()
	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})
	sessionService := NewSessionService(ctx, logger, mockConfig, mockKVRepository)
	timestamp, err := sessionService.GetPodScheduleTimeStamp(mockSessionId)
	assert.Nil(t, err, "sessionService.GetPodScheduleTimeStamp should not throw error")
	assert.Equal(t, mockPodScheduleTimestamp, timestamp)
	mockConfig.AssertNumberOfCalls(t, "Get", 1)
	mockConfig.AssertCalled(t, "Get", "PodScheduleTimeStamp.sessionId")
}

func TestGetPodScheduleTimeStampKeyNotFound(t *testing.T) {
	mockKVRepository := &repository.MockIKVRepository{}
	mockSessionId := "sessionId"
	ctx := context.TODO()
	mockConfig := &config.MockIConfig{}
	mockConfig.On("Get", mock.Anything).Return(nil).Once()
	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})
	sessionService := NewSessionService(ctx, logger, mockConfig, mockKVRepository)
	_, err := sessionService.GetPodScheduleTimeStamp(mockSessionId)
	assert.True(t, NewInvalidStoreKeyErr("PodScheduleTimeStamp.sessionId").Is(err))
}
