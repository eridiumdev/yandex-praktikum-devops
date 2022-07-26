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

func (s *metricsService) Update(metric domain.Metric) (updated domain.Metric, changed bool) {
	existingMetric, found := s.repo.Get(metric.Name)
	if found && metric.IsCounter() {
		// For counters, old value is added on top of new value
		metric.Counter += existingMetric.Counter
	}
	s.repo.Store(metric)
	return metric, metric != existingMetric
}

func (s *metricsService) Get(name string) (metric domain.Metric, found bool) {
	metric, found = s.repo.Get(name)
	return
}

func (s *metricsService) List() []domain.Metric {
	return s.repo.List()
}
