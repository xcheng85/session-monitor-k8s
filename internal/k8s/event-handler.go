package k8s

import "k8s.io/client-go/tools/cache"

//go:generate mockery --name IK8sEventHandler
type IK8sEventHandler interface {
	CustomWatchErrorHandler(r *cache.Reflector, err error)
	OnAddObject(obj interface{})
	OnUpdateObject(oldObj, newObj interface{})
	OnDeleteObject(obj interface{})
}


