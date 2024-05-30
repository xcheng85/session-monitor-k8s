package repository

import (
	"context"
	"encoding/json"
	"testing"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/xcheng85/session-monitor-k8s/internal/config"
	"github.com/xcheng85/session-monitor-k8s/internal/logger"
)

func TestRedisRepository(t *testing.T) {
	ctx := context.TODO()
	mockConfig := &config.MockIConfig{}
	mockConfig.On("Get", "app.redis_disable").Return(false, nil).Once()
	mockConfig.On("Get", "app.redis_address").Return("127.0.0.1:6379", nil).Once()
	mockConfig.On("Get", "app.redis_mock").Return(true, nil).Once()
	mockConfig.On("Get", "app.redisv9_disable").Return(false, nil).Once()
	mockConfig.On("Get", "app.redisv9_address").Return("127.0.0.1:6380", nil).Once()

	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})
	redisRepo, err := NewRedisRepository(ctx, mockConfig, logger)
	assert.NotNil(t, redisRepo)
	assert.Nil(t, err)

	serverTimestamp, err := redisRepo.GetServerTimestamp(ctx)
	assert.NotNil(t, serverTimestamp)
	assert.Nil(t, err)

	resp, err := redisRepo.Ping(ctx)
	assert.Equal(t, "PONG", resp)
	assert.Nil(t, err)
}

func TestRedisRepositoryV9(t *testing.T) {
	ctx := context.TODO()
	mockConfig := &config.MockIConfig{}
	mockConfig.On("Get", "app.redis_disable").Return(false, nil).Once()
	mockConfig.On("Get", "app.redis_address").Return("127.0.0.1:6379", nil).Once()
	mockConfig.On("Get", "app.redis_mock").Return(true, nil).Once()
	mockConfig.On("Get", "app.redisv9_disable").Return(false, nil).Once()
	mockConfig.On("Get", "app.redisv9_address").Return("127.0.0.1:6380", nil).Once()

	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})
	redisRepo, err := NewRedisRepository(ctx, mockConfig, logger)
	assert.NotNil(t, redisRepo)
	assert.Nil(t, err)

	serverTimestamp, err := redisRepo.GetServerTimestamp(ctx)
	assert.NotNil(t, serverTimestamp)
	assert.Nil(t, err)

	resp, err := redisRepo.Ping(ctx)
	assert.Equal(t, "PONG", resp)
	assert.Nil(t, err)
}

func TestAddStreamEvent(t *testing.T) {
	ctx := context.TODO()
	mockKVRepository := &MockIKVRepository{}
	mockServerTimestamp := int64(88888888888)
	mockKVRepository.On("GetServerTimestamp", ctx).Return(mockServerTimestamp, nil).Once()
	mockKVRepository.On("AddStreamEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("taskid-001", nil).Once()
	// 2nd redis client
	mockKVRepository2 := &MockIKVRepository{}
	mockKVRepository2.On("GetServerTimestamp", ctx).Return(mockServerTimestamp, nil).Once()
	mockKVRepository2.On("AddStreamEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("taskid-002", nil).Once()

	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})
	taskInfo := map[string]string{
		"SessionId": "sessionId",
		"NodeName":  "nodeName",
	}
	taskInfoBuf, _ := json.Marshal(taskInfo)
	payloadToSubmit := []interface{}{
		"TaskType",
		"EnqueueSession",
		"TaskInfo",
		string(taskInfoBuf),
		"TaskCreateTimeStamp",
		int64(62135596800),
	}
	redisRepo := &redisRepository{
		clients: []IKVRepository{mockKVRepository, mockKVRepository2},
		logger:  logger,
	}
	res, err := redisRepo.AddStreamEvent(ctx, "enqueue_session_test", "*", payloadToSubmit)
	assert.Equal(t, "taskid-002", res)
	assert.Nil(t, err)
}

func TestAddToUnsortedSet(t *testing.T) {
	ctx := context.TODO()
	mockKVRepository := &MockIKVRepository{}
	mockKVRepository.On("AddToUnsortedSet", mock.Anything, mock.Anything, mock.Anything).Return(int64(1), nil).Once()

	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})

	label := map[string]string{
		"name": "node-1",
	}
	labelBuf, _ := json.Marshal(label)

	redisRepo := &redisRepository{
		clients: []IKVRepository{mockKVRepository},
		logger:  logger,
	}

	res, err := redisRepo.AddToUnsortedSet(ctx, "gpuAgentPoolSetKey", &Object{
		Key:     "viz1",
		Payload: string(labelBuf),
	})
	assert.Equal(t, int64(1), res)
	assert.Nil(t, err)
}
