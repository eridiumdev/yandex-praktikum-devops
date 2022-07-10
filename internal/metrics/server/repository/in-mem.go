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

func (r *inMemRepo) Store(metric domain.Metric) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.metrics[metric.Name()]; ok {
		r.metrics[metric.Name()].Update(metric.Value())
	} else {
		r.metrics[metric.Name()] = metric.Copy()
	}
	return nil
}

func (r *inMemRepo) Get(metricName string) (domain.Metric, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return r.metrics[metricName], nil
}

func (r *inMemRepo) List() ([]domain.Metric, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	result := make([]domain.Metric, 0)
	for _, metric := range r.metrics {
		result = append(result, metric.Copy())
	}
	return result, nil
}
