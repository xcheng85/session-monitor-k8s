package k8s

import (
	"context"
	"errors"

	"github.com/xcheng85/session-monitor-k8s/internal/config"
	"go.uber.org/zap"
	// v1 "k8s.io/api/core/v1"
	// "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	// "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

type k8sDynamicInformer struct {
	logger    *zap.Logger
	config    config.IConfig
	informer  cache.SharedIndexInformer
	ctx       context.Context
	handler   IK8sEventHandler
}

func (informer *k8sDynamicInformer) Run() {
	informer.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    informer.handler.OnAddObject,
		UpdateFunc: informer.handler.OnUpdateObject,
		DeleteFunc: informer.handler.OnDeleteObject,
	})
	informer.informer.SetWatchErrorHandler(informer.handler.CustomWatchErrorHandler)
	informer.informer.Run(informer.ctx.Done())
}

func NewK8sDynamicInformer(ctx context.Context, logger *zap.Logger, config config.IConfig, handler IK8sEventHandler) (IK8sInformer, error) {
	informer, err := newDynamicInformer(ctx, config)
	if err != nil {
		return nil, err
	}

	return &k8sDynamicInformer{
		logger,
		config,
		informer,
		ctx,
		handler,
	}, nil
}

func newDynamicInformer(ctx context.Context, config config.IConfig) (cache.SharedIndexInformer, error) {
	kubeConfig := config.Get("app.kube_config").(string)
	resource := "pods"
	namespace := config.Get("app.pod_namespace").(string)
	var clusterConfig *rest.Config
	var err error
	if kubeConfig != "" {
		clusterConfig, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
	} else {
		clusterConfig, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, err
	}
	dynamicClient, err := dynamic.NewForConfig(clusterConfig)
	if err != nil {
		return nil, err
	}

	podResources := schema.GroupVersionResource{Group: "", Version: "v1", Resource: resource}
	if namespace == "" {
		return nil, errors.New("empty namespace")
	}
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(dynamicClient, 0, namespace, nil)
	informer := factory.ForResource(podResources).Informer()
	return informer, nil
}
