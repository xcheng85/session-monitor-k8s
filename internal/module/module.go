package module

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/xcheng85/session-monitor-k8s/internal/config"
	"go.uber.org/zap"
)

// chi.Mux is the implementation of chi.Router interface
//
//go:generate mockery --name IModuleContext
type IModuleContext interface {
	Mux() *chi.Mux
	Logger() *zap.Logger
	Config() config.IConfig
}

type Module interface {
	Startup(context.Context, IModuleContext) error
}
