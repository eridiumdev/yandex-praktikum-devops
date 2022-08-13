package repository

import (
	"context"
	"sync"

	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

type inMemRepo struct {
	metrics map[string]domain.Metric
	mutex   *sync.RWMutex
}

func NewInMemRepo() *inMemRepo {
	return &inMemRepo{
		metrics: make(map[string]domain.Metric),
		mutex:   &sync.RWMutex{},
	}
}

func (r *inMemRepo) Store(ctx context.Context, metric domain.Metric) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.metrics[metric.Name] = metric
	return nil
}

func (r *inMemRepo) Get(ctx context.Context, name string) (domain.Metric, bool, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	metric, ok := r.metrics[name]
	return metric, ok, nil
}

func (r *inMemRepo) List(ctx context.Context) ([]domain.Metric, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	result := make([]domain.Metric, 0)
	for _, metric := range r.metrics {
		result = append(result, metric)
	}
	return result, nil
}
