package node

import (
	"context"
	
	"github.com/go-chi/chi/v5"
	"github.com/xcheng85/session-monitor-k8s/internal/config"
	"github.com/xcheng85/session-monitor-k8s/internal/ddd"
	"github.com/xcheng85/session-monitor-k8s/internal/k8s"
	"github.com/xcheng85/session-monitor-k8s/internal/module"
	"github.com/xcheng85/session-monitor-k8s/internal/repository"
	"github.com/xcheng85/session-monitor-k8s/internal/session"
	"github.com/xcheng85/session-monitor-k8s/node/internal/handler"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type NodeMonitoringModule struct{}

func (m NodeMonitoringModule) Startup(ctx context.Context, mono module.IModuleContext) (*dig.Container, error) {
	container := dig.New()
	err := container.Provide(func() context.Context {
		return ctx
	})
	err = container.Provide(func() *zap.Logger {
		return mono.Logger()
	})
	err = container.Provide(func() config.IConfig {
		return mono.Config()
	})
	err = container.Provide(func() *chi.Mux {
		return mono.Mux()
	})
	err = container.Provide(func() repository.IKVRepository {
		return mono.KvRepository()
	})
	err = container.Provide(func() ddd.IEventDispatcher[ddd.IEvent] {
		return mono.EventDispatcher()
	})
	err = container.Provide(k8s.NewK8sDynamicInformer)
	if err != nil {
		return nil, err
	}
	err = container.Provide(handler.NewDomainEventHandlers)
	if err != nil {
		return nil, err
	}
	err = container.Provide(handler.NewNodeEventHandler)
	if err != nil {
		return nil, err
	}
	err = container.Provide(func() string {
		return "nodes"
	}, dig.Name("k8s_resource"))
	if err != nil {
		return nil, err
	}
	err = container.Provide(func() string {
		return ""
	}, dig.Name("k8s_resource_namespace"))
	if err != nil {
		return nil, err
	}
	err = container.Provide(session.NewSessionService)
	if err != nil {
		return nil, err
	}
	err = container.Invoke(func(informer k8s.IK8sInformer) error {
		// detach goroutine, let app in the cli to do it
		// go informer.Run()
		return nil
	})
	return container, err
}
func NewNodeMonitoringModule() module.Module {
	return &NodeMonitoringModule{}
}
