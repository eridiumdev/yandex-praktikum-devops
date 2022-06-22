package storage

import (
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics"
	"github.com/pkg/errors"
	"sync"
)

type InMemStorage struct {
	Metrics      map[string]metrics.Metric
	metricsMutex *sync.RWMutex
}

func NewInMemStorage() *InMemStorage {
	return &InMemStorage{
		Metrics:      make(map[string]metrics.Metric),
		metricsMutex: &sync.RWMutex{},
	}
}

func (st *InMemStorage) StoreMetric(metric metrics.Metric) error {
	switch metric.GetType() {
	case metrics.TypeCounter:
		st.metricsMutex.Lock()
		defer st.metricsMutex.Unlock()
		if m, ok := st.Metrics[metric.GetName()]; ok {
			// Increment old metric
			sum := m.GetValue().(metrics.Counter) + metric.GetValue().(metrics.Counter)
			st.Metrics[metric.GetName()] = metrics.CounterMetric{
				AbstractMetric: metrics.AbstractMetric{
					Name: m.GetName(),
				},
				Value: sum,
			}
		} else {
			st.Metrics[metric.GetName()] = metric
		}
	case metrics.TypeGauge:
		fallthrough
	default:
		// Override old metric
		st.Metrics[metric.GetName()] = metric
	}
	return nil
}

func (st *InMemStorage) GetMetric(metricType, metricName string) (metrics.Metric, error) {
	metric, ok := st.Metrics[metricName]
	if !ok {
		return nil, errors.Wrapf(ErrMetricNotFound, "'%s'", metricName)
	}
	if metric.GetType() != metricType {
		return nil, errors.Wrapf(ErrMetricIncorrectType, "'%s'", metricType)
	}
	return metric, nil
}

func (st *InMemStorage) ListMetrics() ([]metrics.Metric, error) {
	result := make([]metrics.Metric, 0)
	for _, m := range st.Metrics {
		result = append(result, m)
	}
	return result, nil
}
