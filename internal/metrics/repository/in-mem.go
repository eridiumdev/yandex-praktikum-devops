package repository

import (
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

func (r *inMemRepo) Store(metric domain.Metric) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.metrics[metric.Name] = metric
}

func (r *inMemRepo) Get(name string) (domain.Metric, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	metric, ok := r.metrics[name]
	return metric, ok
}

func (r *inMemRepo) List() []domain.Metric {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	result := make([]domain.Metric, 0)
	for _, metric := range r.metrics {
		result = append(result, metric)
	}
	return result
}
