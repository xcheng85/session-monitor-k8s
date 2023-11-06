package worker

import (
	"context"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
)

// contract for all the functions passed in to the errgroup
// ctx is for simultanous coordinate
type Worker func(ctx context.Context) error

// D principle of SOLID
// interface in local package for clean arch
type IWorkerSyncer interface {
	Add(fns ...Worker)
	Sync() error
	Context() context.Context
	CancelFunc() context.CancelFunc
}

type workerSyncerImp struct {
	ctx     context.Context
	workers []Worker
	cancel  context.CancelFunc
}

// pointer receiver due to write on state
func (s *workerSyncerImp) Add(w ...Worker) {
	s.workers = append(s.workers, w...)
}

func (s workerSyncerImp) Sync() (err error) {
	// create group with context derive
	g, ctx := errgroup.WithContext(s.ctx)

	g.Go(func() error {
		// triggered if ctx is cancelled
		// The derived Context is canceled the first time a function passed to Go
		// returns a non-nil error or the first time Wait returns, whichever occurs first.
		// Listen for the interrupt signal.
		<-ctx.Done()
		s.cancel()
		return nil
	})

	for _, w := range s.workers {
		w := w
		// create a new goroutine
		g.Go(func() error { return w(ctx) })
	}

	return g.Wait()
}

func (w workerSyncerImp) Context() context.Context {
	return w.ctx
}

func (w workerSyncerImp) CancelFunc() context.CancelFunc {
	return w.cancel
}

type workerSyncerConfig struct {
	parentCtx    context.Context // context.Background()
	catchSignals bool
}

// constructor
func NewWorkerSyncer(parentCtx context.Context) IWorkerSyncer {
	config := &workerSyncerConfig{
		parentCtx:    parentCtx,
		catchSignals: true,
	}
	w := &workerSyncerImp{
		workers: []Worker{},
	}
	w.ctx, w.cancel = context.WithCancel(config.parentCtx)
	if config.catchSignals {
		// gracefully shut down
		w.ctx, w.cancel = signal.NotifyContext(w.ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	}
	return w
}
