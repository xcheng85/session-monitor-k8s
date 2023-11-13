package repository

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"github.com/xcheng85/session-monitor-k8s/internal/config"
	"github.com/xcheng85/session-monitor-k8s/internal/logger"
)

func TestRedisRepository(t *testing.T) {
	ctx := context.TODO()
	mockConfig := &config.MockIConfig{}
	mockConfig.On("Get", "app.redis_address").Return("127.0.0.1:6379", nil).Once()
	mockConfig.On("Get", "app.redis_mock").Return(true, nil).Once()

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
	db, mock := redismock.NewClientMock()
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

	mock.ExpectXAdd(&redis.XAddArgs{
		Stream: "enqueue_session_test",
		ID:     "*",
		Values: payloadToSubmit,
	}).SetVal("taskid-001")

	redisRepo := &redisRepository{
		client: db,
		logger: logger,
	}
	res, err := redisRepo.AddStreamEvent(ctx, "enqueue_session_test", "*", payloadToSubmit)
	assert.Equal(t, "taskid-001", res)
	assert.Nil(t, err)
}

func TestAddToUnsortedSet(t *testing.T) {
	ctx := context.TODO()
	db, mock := redismock.NewClientMock()
	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})

	label := map[string]string{
		"name": "node-1",
	}
	labelBuf, _ := json.Marshal(label)

	mock.ExpectTxPipeline()
	mock.ExpectSet("viz1", string(labelBuf), 0).SetVal("OK")
	mock.ExpectSAdd("gpuAgentPoolSetKey", "viz1").SetVal(1)
	mock.ExpectTxPipelineExec()

	redisRepo := &redisRepository{
		client: db,
		logger: logger,
	}

	res, err := redisRepo.AddToUnsortedSet(ctx, "gpuAgentPoolSetKey", &Object{
		Key:     "viz1",
		Payload: string(labelBuf),
	})
	assert.Equal(t, int64(1), res)
	assert.Nil(t, err)
}
