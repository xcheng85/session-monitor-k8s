package k8s

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/xcheng85/session-monitor-k8s/internal/module"
	"testing"
	"time"
)

func Test_ModuleStartup(t *testing.T) {
	mux := chi.NewRouter()
	mockModuleCtx := module.NewMockIModuleContext(t)
	mockModuleCtx.On("Mux").Return(mux).Once()
	module := NewK8sModule()
	// define context and therefore test timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	error := module.Startup(ctx, mockModuleCtx)
	assert.Nil(t, error, "k8s module can start up")
}
