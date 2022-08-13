package service

import (
	"context"

	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

// These are the interfaces required for the Service to work

// MetricsRepository should store and retrieve metrics using backend storage
type MetricsRepository interface {
	Store(ctx context.Context, metric domain.Metric) error
	Get(ctx context.Context, name string) (m domain.Metric, found bool, err error)
	List(ctx context.Context) ([]domain.Metric, error)
}

// MetricsBackuper should be able to backup and restore metrics using long-term storage
type MetricsBackuper interface {
	Backup(metrics []domain.Metric) error
	Restore() ([]domain.Metric, error)
}
