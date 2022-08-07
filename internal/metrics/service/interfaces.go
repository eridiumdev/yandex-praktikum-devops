package service

import (
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

// These are the interfaces required for the Service to work

// MetricsRepository should store and retrieve metrics using backend storage
type MetricsRepository interface {
	Store(metric domain.Metric)
	Get(name string) (metric domain.Metric, found bool)
	List() []domain.Metric
}

// MetricsBackuper should be able to backup and restore metrics using long-term storage
type MetricsBackuper interface {
	Backup(metrics []domain.Metric) error
	Restore() ([]domain.Metric, error)
}
