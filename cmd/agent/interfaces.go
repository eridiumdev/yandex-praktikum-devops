package main

import (
	"context"

	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

// These are the interfaces required by the Agent to work with metrics

// Executor bundles together common functionality for collectors/exporters
// Executor has a name and can become idle/ready for execution
type Executor interface {
	Name() string
	Ready() <-chan bool
}

// MetricsCollector is an Executor that should collect metrics from somewhere
type MetricsCollector interface {
	Executor
	Collect(context.Context) ([]domain.Metric, error)
}

// MetricsExporter is an Executor that should export metrics to somewhere
type MetricsExporter interface {
	Executor
	Export(context.Context, []domain.Metric) error
}

// MetricsBufferer should buffer metrics in temporary storage before exporting
type MetricsBufferer interface {
	Buffer([]domain.Metric)
	Retrieve() []domain.Metric
	Flush()
}
