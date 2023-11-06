package handler

import (
	"github.com/xcheng85/session-monitor-k8s/internal/k8s"
	"go.uber.org/zap"
	"k8s.io/client-go/tools/cache"
)

type PodEventHandler struct {
	logger *zap.Logger
}

func NewPodEventHandler(logger *zap.Logger) k8s.IK8sEventHandler {
	return &PodEventHandler{
		logger,
	}
}

func (handler *PodEventHandler) CustomWatchErrorHandler(r *cache.Reflector, err error) {
	handler.logger.Sugar().Errorw("Watch error", err)
}

func (handler *PodEventHandler) OnAddObject(obj interface{}) {
	handler.logger.Sugar().Infow("onAddObject", obj)
}

func (handler *PodEventHandler) OnUpdateObject(oldObj, newObj interface{}) {
	handler.logger.Sugar().Infow("onUpdateObject", oldObj, newObj)
}

func (handler *PodEventHandler) OnDeleteObject(obj interface{}) {
	handler.logger.Sugar().Infow("onDeleteObject", obj)
}
