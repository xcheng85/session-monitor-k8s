package main

import (
	"fmt"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/xcheng85/session-monitor-k8s/internal/config"
	"github.com/xcheng85/session-monitor-k8s/internal/logger"
	"github.com/xcheng85/session-monitor-k8s/internal/module"
	"github.com/xcheng85/session-monitor-k8s/internal/worker"
	"github.com/xcheng85/session-monitor-k8s/k8s"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type IocContainer struct {
	container *dig.Container
}

func newIocContainer() (*IocContainer, error) {
	container := dig.New()
	container.Provide(newContext)
	err := container.Provide(
		func() *zap.Logger {
			return logger.NewZapLogger(logger.LogConfig{
				LogLevel: logger.DEBUG,
			})
		})
	if err != nil {
		return nil, err
	}
	err = container.Provide(
		func(logger *zap.Logger) (config.IConfig, error) {
			return config.NewViperConfig("./dummy.yaml", []string{os.Getenv("CONFIG_PATH")}, logger)
		})
	err = container.Provide(k8s.NewK8sModule, dig.Name("k8s"))
	err = container.Provide(newMux)
	err = container.Provide(newModuleContext)
	err = container.Provide(worker.NewWorkerSyncer)
	err = container.Provide(func(p struct {
		dig.In
		ModuleContext module.IModuleContext
		K8s           module.Module `name:"k8s"`
		Mux           *chi.Mux
		WorkerSyncer  worker.IWorkerSyncer
	}) *CompositionRoot {
		root := newCompositionRoot(p.Mux, p.ModuleContext, p.K8s, p.WorkerSyncer)
		root.startupModules()
		return root
	})
	return &IocContainer{
		container,
	}, err
}

func (ioc *IocContainer) start() (err error) {
	err = ioc.container.Invoke(func(root *CompositionRoot) error {
		root.startup()
		return nil
	})
	return err
}

func main() {
	ioc, err := newIocContainer()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	err = ioc.start()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
