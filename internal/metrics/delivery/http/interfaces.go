package http

import (
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

// These are the interfaces required for handling metrics requests

// MetricsRenderer should apply metrics to some template, resulting in renderable output
type MetricsRenderer interface {
	RenderList(list []domain.Metric) ([]byte, error)
}

// MetricsService should be able to perform common operations on metrics, such as updating and retrieving
type MetricsService interface {
	Update(metric domain.Metric) (updated domain.Metric, changed bool)
	Get(name string) (metric domain.Metric, found bool)
	List() []domain.Metric
}
