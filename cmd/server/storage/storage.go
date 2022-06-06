package storage

import (
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics"
)

type Storage interface {
	StoreMetric(metric metrics.Metric) error
}
