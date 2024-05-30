package node

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/xcheng85/session-monitor-k8s/internal/config"
	"github.com/xcheng85/session-monitor-k8s/internal/ddd"
	"github.com/xcheng85/session-monitor-k8s/internal/logger"
	"github.com/xcheng85/session-monitor-k8s/internal/module"
	"github.com/xcheng85/session-monitor-k8s/internal/repository"
)

func Test_ModuleStartup(t *testing.T) {
	//mux := chi.NewRouter()
	mockConfig := config.NewMockIConfig(t)
	mockModuleCtx := module.NewMockIModuleContext(t)
	mockLogger := logger.NewZapLogger(logger.LogConfig{
		LogLevel: logger.DEBUG,
	})
	mockEventDispatcher := ddd.NewMockIEventDispatcher[ddd.IEvent](t)
	mockKVRepository := &repository.MockIKVRepository{}

	mockModuleCtx.On("Logger").Return(mockLogger).Once()
	mockModuleCtx.On("Config").Return(mockConfig).Once()
	mockModuleCtx.On("EventDispatcher").Return(mockEventDispatcher).Once()
	mockModuleCtx.On("KvRepository").Return(mockKVRepository).Once()
	mockEventDispatcher.On("Subscribe", mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockConfig.On("Get", "app.kube_config").Return("", nil).Once()

	module := NewNodeMonitoringModule()
	// define context and therefore test timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := module.Startup(ctx, mockModuleCtx)
	assert.NotNil(t, err, "node module cannot start up without valid kube_config")
}
