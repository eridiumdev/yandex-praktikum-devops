package worker

import (
	"context"

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

func (w *Worker) Name() string {
	return w.name
}

func (w *Worker) MaxThreads() int {
	return w.maxThreads
}

func (w *Worker) Reserve(ctx context.Context) bool {
	select {
	case <-w.available:
		logger.New(ctx).Debugf("[%s worker] I have been reserved", w.name)
		return true
	case <-ctx.Done():
		logger.New(ctx).Debugf("[%s worker] Context canceled, reserve aborted", w.name)
		return false
	}
}

func (w *Worker) Release(ctx context.Context) bool {
	select {
	case w.available <- struct{}{}:
		logger.New(ctx).Debugf("[%s worker] I have been released", w.name)
		return true
	case <-ctx.Done():
		logger.New(ctx).Debugf("[%s worker] Context canceled, release aborted", w.name)
		return false
	}
}
