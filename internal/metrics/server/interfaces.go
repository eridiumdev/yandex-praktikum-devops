package server

import (
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

// These are the interfaces required by the Server to work with metrics

// MetricsRenderer should apply metrics to template file, resulting in renderable output
type MetricsRenderer interface {
	RenderList(templateName string, data []domain.Metric) ([]byte, error)
}

// MetricsRepository should store and retrieve metrics using backend storage
type MetricsRepository interface {
	Store(metric domain.Metric) error
	Get(metricName string) (domain.Metric, error)
	List() ([]domain.Metric, error)
}
