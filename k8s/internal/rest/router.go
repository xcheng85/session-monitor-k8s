package rest

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/xcheng85/session-monitor-k8s/k8s/internal/handler"
)

type K8sRouter struct {
	handler handler.IK8sHandler
	mux     *chi.Mux
	ctx     context.Context
}

func NewK8sRouter(handler handler.IK8sHandler, ctx context.Context, mux *chi.Mux) *K8sRouter {
	return &K8sRouter{
		handler: handler,
		mux:     mux,
		ctx:     ctx,
	}
}

func (router *K8sRouter) Register() error {
	r := chi.NewRouter()
	r.Get("/livenessProbe", router.handler.GetLivenessProbe)
	r.Get("/readinessProbe", router.handler.GetReadinessProbe)
	// mounting path must be unique
	router.mux.Mount("/k8s", r)
	return nil
}
