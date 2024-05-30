package module

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/xcheng85/session-monitor-k8s/internal/config"
	"github.com/xcheng85/session-monitor-k8s/internal/ddd"
	"github.com/xcheng85/session-monitor-k8s/internal/repository"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

// chi.Mux is the implementation of chi.Router interface
// all the singleton/stateful is provided in the ModuleContext
//
//go:generate mockery --name IModuleContext
type IModuleContext interface {
	Mux() *chi.Mux
	Logger() *zap.Logger
	Config() config.IConfig
	KvRepository() repository.IKVRepository
	EventDispatcher() ddd.IEventDispatcher[ddd.IEvent] // translate k8s event into domain event
}

type Module interface {
	Startup(context.Context, IModuleContext) (*dig.Container, error)
}
