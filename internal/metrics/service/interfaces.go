package service

import (
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

// These are the interfaces required for the Service to work

// MetricsRepository should store and retrieve metrics using backend storage
type MetricsRepository interface {
	Store(metric domain.Metric) error
	Get(metricName string) (domain.Metric, error)
	List() ([]domain.Metric, error)
}
