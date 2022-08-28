package worker

import (
	"context"

	"eridiumdev/yandex-praktikum-go-devops/internal/common/helpers"
	"eridiumdev/yandex-praktikum-go-devops/internal/common/logger"
)

type Worker struct {
	name       string
	maxThreads int
	available  chan struct{}
}

func New(name string, maxThreads int) *Worker {
	w := &Worker{
		name:       name,
		maxThreads: maxThreads,
		available:  make(chan struct{}, maxThreads),
	}
	// Make all threads available at the start
	for i := 0; i < maxThreads; i++ {
		w.available <- struct{}{}
	}
	return w
}

func (w *Worker) LogSource() string {
	return w.name + " worker"
}

func (w *Worker) Name() string {
	return w.name
}

func (w *Worker) MaxThreads() int {
	return w.maxThreads
}

func (w *Worker) Reserve(ctx context.Context) error {
	select {
	case <-w.available:
		logger.New(ctx).Src(w).Debugf("I have been reserved")
		return nil
	case <-ctx.Done():
		return helpers.NewErr(w, "still busy (context timeout)")
	}
}

func (w *Worker) Release(ctx context.Context) error {
	select {
	case w.available <- struct{}{}:
		logger.New(ctx).Src(w).Debugf("I have been released")
		return nil
	case <-ctx.Done():
		return helpers.NewErr(w, "all workers already available (chan buffer full)")
	}
}
