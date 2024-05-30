package k8s

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/xcheng85/session-monitor-k8s/internal/module"
	"github.com/xcheng85/session-monitor-k8s/k8s/internal/handler"
	"github.com/xcheng85/session-monitor-k8s/k8s/internal/rest"
	"go.uber.org/dig"
)

type K8sModule struct{}

func (m K8sModule) Startup(ctx context.Context, mono module.IModuleContext) (*dig.Container, error) {
	container := dig.New()
	container.Provide(handler.NewK8sHandler)
	container.Provide(rest.NewK8sRouter)
	container.Provide(func() *chi.Mux {
		return mono.Mux()
	})
	container.Provide(func() context.Context {
		return ctx
	})
	err := container.Invoke(func(r *rest.K8sRouter) error {
		return r.Register()
	})
	return container, err
}

func NewK8sModule() module.Module {
	return &K8sModule{}
}
