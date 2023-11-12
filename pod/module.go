package pod

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/xcheng85/session-monitor-k8s/internal/config"
	"github.com/xcheng85/session-monitor-k8s/internal/ddd"
	"github.com/xcheng85/session-monitor-k8s/internal/k8s"
	"github.com/xcheng85/session-monitor-k8s/internal/module"
	"github.com/xcheng85/session-monitor-k8s/internal/repository"
	"github.com/xcheng85/session-monitor-k8s/internal/session"
	"github.com/xcheng85/session-monitor-k8s/pod/internal/handler"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type PodMonitoringModule struct{}

func (m PodMonitoringModule) Startup(ctx context.Context, mono module.IModuleContext) (err error) {
	container := dig.New()
	err = container.Provide(func() context.Context {
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
		return err
	}
	err = container.Provide(handler.NewDomainEventHandlers)
	if err != nil {
		return err
	}
	err = container.Provide(handler.NewPodEventHandler)
	if err != nil {
		return err
	}
	err = container.Provide(func() string {
		return "pods"
	}, dig.Name("k8s_resource"))
	if err != nil {
		return err
	}
	err = container.Provide(session.NewSessionService)
	if err != nil {
		return err
	}
	err = container.Invoke(func(informer k8s.IK8sInformer) error {
		// detach goroutine
		go informer.Run()
		return nil
	})
	return err
}
func NewPodMonitoringModule() module.Module {
	return &PodMonitoringModule{}
}
