package storage

import (
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics"
	"errors"
)

type Storage interface {
	StoreMetric(metric metrics.Metric) error
	GetMetric(metricType, metricName string) (metrics.Metric, error)
	ListMetrics() ([]metrics.Metric, error)
}

var (
	ErrMetricNotFound      = errors.New("metric not found")
	ErrMetricIncorrectType = errors.New("incorrect metric type")
)
