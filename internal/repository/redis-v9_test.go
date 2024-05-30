package repository

import (
	"context"
	"github.com/xcheng85/session-monitor-k8s/internal/logger"
	"encoding/json"
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRedisClientV9_AddStreamEvent(t *testing.T) {
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

	//v9, err := newRedisClientV9("127.0.0.1:6379", true, logger)
	v9 := &redisClientV9{
		db,
		logger,
	}
	res, err := v9.AddStreamEvent(ctx, "enqueue_session_test", "*", payloadToSubmit)
	assert.Equal(t, "taskid-001", res)
	assert.Nil(t, err)
}

func TestRedisClientV9_AddToUnsortedSet(t *testing.T) {
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

	//v9, err := newRedisClientV9("127.0.0.1:6379", true, logger)
	v9 := &redisClientV9{
		db,
		logger,
	}

	res, err := v9.AddToUnsortedSet(ctx, "gpuAgentPoolSetKey", &Object{
		Key:     "viz1",
		Payload: string(labelBuf),
	})
	assert.Equal(t, int64(1), res)
	assert.Nil(t, err)
}
