package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/xcheng85/session-monitor-k8s/internal/ddd"
	"github.com/xcheng85/session-monitor-k8s/internal/logger"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestNewPodEventHandler(t *testing.T) {
	ctx := context.TODO()
	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})
	eventDispatcher := ddd.MockIEventDispatcher[ddd.IEvent]{}
	eventHandler := ddd.MockIEventHandler[ddd.IEvent]{}

	h := NewPodEventHandler(ctx, logger, &eventDispatcher, &eventHandler)
	assert.NotNil(t, h, "Pod Event Handler should not be nil")
}

func TestCustomWatchErrorHandler(t *testing.T) {
	ctx := context.TODO()
	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})
	eventDispatcher := ddd.MockIEventDispatcher[ddd.IEvent]{}
	eventHandler := ddd.MockIEventHandler[ddd.IEvent]{}

	err := errors.New("errors from k8s api server")
	eventDispatcher.On("Publish", ctx, mock.Anything).Return(nil).Once()

	h := NewPodEventHandler(ctx, logger, &eventDispatcher, &eventHandler)
	h.CustomWatchErrorHandler(nil, err)

	eventDispatcher.AssertNumberOfCalls(t, "Publish", 1)
	// 2nd argument is dynamic, use anything.
	eventDispatcher.AssertCalled(t, "Publish", ctx, mock.Anything)
}

func TestOnAddObject(t *testing.T) {
	ctx := context.TODO()
	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})

	eventDispatcher := ddd.MockIEventDispatcher[ddd.IEvent]{}
	eventDispatcher.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	eventHandler := ddd.MockIEventHandler[ddd.IEvent]{}
	h := NewPodEventHandler(ctx, logger, &eventDispatcher, &eventHandler)
	payload := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "Pod",
			"apiVersion": "v1",
			"metadata": map[string]interface{}{
				"name":            "test-app-pod",
				"namespace":       "test-namespace",
				"uid":             "test_uid",
				"resourceVersion": "test_resourceVersion",
				"labels": map[string]interface{}{
					"sessionId": "test-app",
				},
				"annotations": map[string]interface{}{
					"nfd.node.kubernetes.io/worker.version": "v0.14.1",
				},
			},
			"spec": map[string]interface{}{},
		},
	}
	h.OnAddObject(payload)
	eventDispatcher.AssertNumberOfCalls(t, "Publish", 1)
}

func TestOnUpdateObject_Failed(t *testing.T) {
	ctx := context.TODO()
	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})

	eventDispatcher := ddd.MockIEventDispatcher[ddd.IEvent]{}
	eventDispatcher.On("Publish", ctx, mock.Anything).Return(nil).Once()

	eventHandler := ddd.MockIEventHandler[ddd.IEvent]{}
	h := NewPodEventHandler(ctx, logger, &eventDispatcher, &eventHandler)
	payload := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "Pod",
			"apiVersion": "v1",
			"metadata": map[string]interface{}{
				"name":      "test_name",
				"namespace": "test_namespace",
				"labels": map[string]interface{}{
					"sessionId": "session-123",
				},
			},
			"spec": map[string]interface{}{
				"nodeName": "node-123",
			},
			"status": map[string]interface{}{
				"phase": "Failed",
			},
		},
	}
	h.OnUpdateObject(nil, payload)
	eventDispatcher.AssertNumberOfCalls(t, "Publish", 1)
}

func TestOnUpdateObject_Succeeded(t *testing.T) {
	ctx := context.TODO()
	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})

	eventDispatcher := ddd.MockIEventDispatcher[ddd.IEvent]{}
	eventDispatcher.On("Publish", ctx, mock.Anything).Return(nil).Once()

	eventHandler := ddd.MockIEventHandler[ddd.IEvent]{}
	h := NewPodEventHandler(ctx, logger, &eventDispatcher, &eventHandler)
	payload := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "Pod",
			"apiVersion": "v1",
			"metadata": map[string]interface{}{
				"name":      "test_name",
				"namespace": "test_namespace",
				"labels": map[string]interface{}{
					"sessionId": "session-123",
				},
			},
			"spec": map[string]interface{}{
				"nodeName": "node-123",
			},
			"status": map[string]interface{}{
				"phase": "Succeeded",
			},
		},
	}
	h.OnUpdateObject(nil, payload)
	eventDispatcher.AssertNumberOfCalls(t, "Publish", 1)
}

func TestOnUpdateObject_Unknown(t *testing.T) {
	ctx := context.TODO()
	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})

	eventDispatcher := ddd.MockIEventDispatcher[ddd.IEvent]{}
	eventDispatcher.On("Publish", ctx, mock.Anything).Return(nil).Once()

	eventHandler := ddd.MockIEventHandler[ddd.IEvent]{}
	h := NewPodEventHandler(ctx, logger, &eventDispatcher, &eventHandler)
	payload := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "Pod",
			"apiVersion": "v1",
			"metadata": map[string]interface{}{
				"name":      "test_name",
				"namespace": "test_namespace",
				"labels": map[string]interface{}{
					"sessionId": "session-123",
				},
			},
			"spec": map[string]interface{}{
				"nodeName": "node-123",
			},
			"status": map[string]interface{}{
				"phase": "Unknown",
			},
		},
	}
	h.OnUpdateObject(nil, payload)
	eventDispatcher.AssertNumberOfCalls(t, "Publish", 1)
}

func TestOnUpdateObject_Running(t *testing.T) {
	ctx := context.TODO()
	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})

	eventDispatcher := ddd.MockIEventDispatcher[ddd.IEvent]{}
	eventDispatcher.On("Publish", ctx, mock.Anything).Return(nil).Once()

	eventHandler := ddd.MockIEventHandler[ddd.IEvent]{}
	h := NewPodEventHandler(ctx, logger, &eventDispatcher, &eventHandler)
	payload := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "Pod",
			"apiVersion": "v1",
			"metadata": map[string]interface{}{
				"name":      "test_name",
				"namespace": "test_namespace",
				"labels": map[string]interface{}{
					"sessionId": "session-123",
				},
			},
			"spec": map[string]interface{}{
				"nodeName": "node-123",
			},
			"status": map[string]interface{}{
				"podIP": "1.2.3.4",
				"phase": "Running",
				"conditions": []map[string]interface{}{
					{
						"type":   "Initialized",
						"status": "True",
					},
					{
						"type":   "Ready",
						"status": "True",
					},
					{
						"type":   "ContainersReady",
						"status": "True",
					},
					{
						"type":   "PodScheduled",
						"status": "True",
					},
				},
			},
		},
	}
	h.OnUpdateObject(nil, payload)
	eventDispatcher.AssertNumberOfCalls(t, "Publish", 1)
}

// AKS pod in unstable state between deleted and running
func TestOnUpdateObject_Running_Unstable(t *testing.T) {
	ctx := context.TODO()
	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})

	eventDispatcher := ddd.MockIEventDispatcher[ddd.IEvent]{}
	eventDispatcher.On("Publish", ctx, mock.Anything).Return(nil).Once()

	eventHandler := ddd.MockIEventHandler[ddd.IEvent]{}
	h := NewPodEventHandler(ctx, logger, &eventDispatcher, &eventHandler)
	payload := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "Pod",
			"apiVersion": "v1",
			"metadata": map[string]interface{}{
				"name":      "test_name",
				"namespace": "test_namespace",
				"labels": map[string]interface{}{
					"sessionId": "session-123",
				},
				"deletionTimestamp": "2023-11-10T23:00:00Z",
			},
			"spec": map[string]interface{}{
				"nodeName": "node-123",
			},
			"status": map[string]interface{}{
				"podIP": "1.2.3.4",
				"phase": "Running",
				"conditions": []map[string]interface{}{
					{
						"type":   "Initialized",
						"status": "True",
					},
					{
						"type":   "Ready",
						"status": "True",
					},
					{
						"type":   "ContainersReady",
						"status": "True",
					},
					{
						"type":   "PodScheduled",
						"status": "True",
					},
				},
			},
		},
	}
	h.OnUpdateObject(nil, payload)
	eventDispatcher.AssertNumberOfCalls(t, "Publish", 1)
}

func TestOnUpdateObject_Pending(t *testing.T) {
	ctx := context.TODO()
	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})

	eventDispatcher := ddd.MockIEventDispatcher[ddd.IEvent]{}
	eventDispatcher.On("Publish", ctx, mock.Anything).Return(nil).Once()

	eventHandler := ddd.MockIEventHandler[ddd.IEvent]{}
	h := NewPodEventHandler(ctx, logger, &eventDispatcher, &eventHandler)
	payload := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "Pod",
			"apiVersion": "v1",
			"metadata": map[string]interface{}{
				"name":      "test_name",
				"namespace": "test_namespace",
				"labels": map[string]interface{}{
					"sessionId": "session-123",
				},
			},
			"spec": map[string]interface{}{
				"nodeName": "node-123",
			},
			"status": map[string]interface{}{
				"phase": "Pending",
				"conditions": []map[string]interface{}{
					{
						"type":   "Initialized",
						"status": "True",
					},
					{
						"type":   "PodScheduled",
						"status": "True",
					},
				},
			},
		},
	}
	h.OnUpdateObject(nil, payload)
	eventDispatcher.AssertNumberOfCalls(t, "Publish", 1)
}

func TestOnUpdateObject_Crash(t *testing.T) {
	ctx := context.TODO()
	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})

	eventDispatcher := ddd.MockIEventDispatcher[ddd.IEvent]{}
	eventDispatcher.On("Publish", ctx, mock.Anything).Return(nil).Once()

	eventHandler := ddd.MockIEventHandler[ddd.IEvent]{}
	h := NewPodEventHandler(ctx, logger, &eventDispatcher, &eventHandler)
	payload := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "Pod",
			"apiVersion": "v1",
			"metadata": map[string]interface{}{
				"name":      "test_name",
				"namespace": "test_namespace",
				"labels": map[string]interface{}{
					"sessionId": "session-123",
				},
			},
			"spec": map[string]interface{}{
				"nodeName": "node-123",
			},
			"status": map[string]interface{}{
				"phase": "Running",
				"conditions": []map[string]interface{}{
					{
						"type":   "Initialized",
						"status": "True",
					},
				},
				"containerStatuses": []map[string]interface{}{
					{
						"name": "container-0",
						"state": map[string]interface{}{
							"terminated": map[string]interface{}{
								"startedAt": "2022-01-01T15:04:05Z",
							},
						},
					},
				},
			},
		},
	}
	h.OnUpdateObject(nil, payload)
	eventDispatcher.AssertNumberOfCalls(t, "Publish", 1)
}

func TestOnDeleteObject(t *testing.T) {
	ctx := context.TODO()
	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})

	eventDispatcher := ddd.MockIEventDispatcher[ddd.IEvent]{}
	eventDispatcher.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	eventHandler := ddd.MockIEventHandler[ddd.IEvent]{}
	h := NewPodEventHandler(ctx, logger, &eventDispatcher, &eventHandler)
	payload := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "Pod",
			"apiVersion": "v1",
			"metadata": map[string]interface{}{
				"name":            "test-app-pod",
				"namespace":       "test-namespace",
				"uid":             "test_uid",
				"resourceVersion": "test_resourceVersion",
				"labels": map[string]interface{}{
					"sessionId": "test-app",
				},
				"annotations": map[string]interface{}{
					"nfd.node.kubernetes.io/worker.version": "v0.14.1",
				},
			},
			"spec": map[string]interface{}{},
		},
	}
	h.OnDeleteObject(payload)
	eventDispatcher.AssertNumberOfCalls(t, "Publish", 1)
}
