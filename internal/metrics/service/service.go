package service

import (
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

type metricsService struct {
	repo MetricsRepository
}

func NewMetricsService(repo MetricsRepository) *metricsService {
	return &metricsService{
		repo: repo,
	}
}

func (s *metricsService) UpdateCounter(name string, value domain.Counter) (domain.Metric, error) {
	metric, err := s.repo.Get(name)
	if err != nil {
		return nil, err
	}
	if metric == nil {
		metric = domain.NewCounter(name, value)
	} else {
		// For counters, new value is added on top of previous value
		err = metric.Add(value)
		if err != nil {
			return metric, err
		}
	}
	return metric, s.repo.Store(metric)
}

func (s *metricsService) UpdateGauge(name string, value domain.Gauge) (domain.Metric, error) {
	metric, err := s.repo.Get(name)
	if err != nil {
		return nil, err
	}
	if metric == nil {
		metric = domain.NewGauge(name, value)
	} else {
		// For gauges, previous value is overwritten
		err = metric.Set(value)
		if err != nil {
			return metric, err
		}
	}
	return metric, s.repo.Store(metric)
}

func (s *metricsService) Get(name string) (domain.Metric, error) {
	return s.repo.Get(name)
}

func (s *metricsService) List() ([]domain.Metric, error) {
	return s.repo.List()
}
