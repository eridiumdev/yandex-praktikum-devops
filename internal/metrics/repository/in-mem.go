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

func (r *inMemRepo) Store(ctx context.Context, metrics ...domain.Metric) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, metric := range metrics {
		r.metrics[metric.Name] = metric
	}
	return nil
}

func (r *inMemRepo) Get(ctx context.Context, name string) (domain.Metric, bool, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	metric, ok := r.metrics[name]
	return metric, ok, nil
}

func (r *inMemRepo) List(ctx context.Context, filter *domain.MetricsFilter) ([]domain.Metric, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	result := make([]domain.Metric, 0)

	for _, metric := range r.metrics {
		if filter != nil && len(filter.Names) > 0 && !sliceContains(filter.Names, metric.Name) {
			// Skip if metric name not in filter
			continue
		}
		result = append(result, metric)
	}
	return result, nil
}

func sliceContains(slice []string, elem string) bool {
	for _, value := range slice {
		if value == elem {
			return true
		}
	}
	return false
}
