package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/xcheng85/session-monitor-k8s/internal/config"
	"github.com/xcheng85/session-monitor-k8s/internal/ddd"
	"github.com/xcheng85/session-monitor-k8s/internal/logger"
)

func TestNewNodeEventHandler(t *testing.T) {
	ctx := context.TODO()
	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})
	config := &config.MockIConfig{}
	eventDispatcher := ddd.MockIEventDispatcher[ddd.IEvent]{}
	eventHandler := ddd.MockIEventHandler[ddd.IEvent]{}

	h := NewNodeEventHandler(ctx, logger, config, &eventDispatcher, &eventHandler)
	assert.NotNil(t, h, "Node Event Handler should not be nil")
}

func TestCustomWatchErrorHandler(t *testing.T) {
	ctx := context.TODO()
	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})
	config := &config.MockIConfig{}
	eventDispatcher := ddd.MockIEventDispatcher[ddd.IEvent]{}
	eventHandler := ddd.MockIEventHandler[ddd.IEvent]{}

	err := errors.New("errors from k8s api server")
	eventDispatcher.On("Publish", ctx, mock.Anything).Return(nil).Once()

	h := NewNodeEventHandler(ctx, logger, config, &eventDispatcher, &eventHandler)
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

	config := &config.MockIConfig{}
	args := []interface{}{"viz3d", "viz3d1"}
	config.On("Get", "app.gpu_agent_pool_names").Return(args).Once()

	eventDispatcher := ddd.MockIEventDispatcher[ddd.IEvent]{}
	eventDispatcher.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	eventHandler := ddd.MockIEventHandler[ddd.IEvent]{}
	h := NewNodeEventHandler(ctx, logger, config, &eventDispatcher, &eventHandler)
	payload := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "Node",
			"apiVersion": "v1",
			"metadata": map[string]interface{}{
				"name":            "aks-viz3d4-33002848-vmss0001nc",
				"uid":             "test_uid",
				"resourceVersion": "test_resourceVersion",
				"labels": map[string]interface{}{
					"accelerator": "nvidia",
					"agentpool":   "viz3d",
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

func TestOnUpdateObject(t *testing.T) {
	ctx := context.TODO()
	logger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})

	config := &config.MockIConfig{}
	args := []interface{}{"viz3d", "viz3d1"}
	config.On("Get", "app.gpu_agent_pool_names").Return(args).Once()

	eventDispatcher := ddd.MockIEventDispatcher[ddd.IEvent]{}
	eventDispatcher.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	eventHandler := ddd.MockIEventHandler[ddd.IEvent]{}
	h := NewNodeEventHandler(ctx, logger, config, &eventDispatcher, &eventHandler)
	payload := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "Node",
			"apiVersion": "v1",
			"metadata": map[string]interface{}{
				"name":            "aks-viz3d4-33002848-vmss0001nc",
				"uid":             "test_uid",
				"resourceVersion": "test_resourceVersion",
				"labels": map[string]interface{}{
					"accelerator": "nvidia",
					"agentpool":   "viz3d",
				},
				"annotations": map[string]interface{}{
					"nfd.node.kubernetes.io/worker.version": "v0.14.1",
				},
			},
			"spec": map[string]interface{}{},
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

	config := &config.MockIConfig{}
	args := []interface{}{"viz3d", "viz3d1"}
	config.On("Get", "app.gpu_agent_pool_names").Return(args).Once()

	eventDispatcher := ddd.MockIEventDispatcher[ddd.IEvent]{}
	eventDispatcher.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	eventHandler := ddd.MockIEventHandler[ddd.IEvent]{}
	h := NewNodeEventHandler(ctx, logger, config, &eventDispatcher, &eventHandler)
	payload := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "Node",
			"apiVersion": "v1",
			"metadata": map[string]interface{}{
				"name":            "aks-viz3d4-33002848-vmss0001nc",
				"uid":             "test_uid",
				"resourceVersion": "test_resourceVersion",
				"labels": map[string]interface{}{
					"accelerator": "nvidia",
					"agentpool":   "viz3d",
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

func TestParseNode(t *testing.T) {
	payload := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "Node",
			"apiVersion": "v1",
			"metadata": map[string]interface{}{
				"name":            "aks-viz3d4-33002848-vmss0001nc",
				"uid":             "test_uid",
				"resourceVersion": "test_resourceVersion",
				"labels": map[string]interface{}{
					"accelerator": "nvidia",
					"agentpool":   "viz3d",
				},
				"annotations": map[string]interface{}{
					"nfd.node.kubernetes.io/worker.version": "v0.14.1",
				},
			},
			"spec": map[string]interface{}{},
		},
	}

	node, err := parseNode(payload)
	assert.NotNil(t, node, "parseNode should return valid node")
	assert.Equal(t, "v1", node.APIVersion)
	assert.Equal(t, "Node", node.Kind)
	assert.Nil(t, err, "parseNode should not throw error")
}

func TestParseGPUDriverVersion(t *testing.T) {
	nodeLables := map[string]string{
		"nvidia.com/cuda.driver.major": "525",
		"nvidia.com/cuda.driver.minor": "85",
		"nvidia.com/cuda.driver.rev":   "12",
	}
	driverVersion := parseGPUDriverVersion(&nodeLables)
	assert.Equal(t, "525.85.12", driverVersion, "Parse GPU Driver Version should pass")
}

func TestParseGPUDriverVersionMissingLables(t *testing.T) {
	nodeLables := map[string]string{
		"nvidia.com/cuda.driver.major": "525",
	}
	driverVersion := parseGPUDriverVersion(&nodeLables)
	assert.Equal(t, "", driverVersion, "Parse invalid gpu labels should return empty string")
}
