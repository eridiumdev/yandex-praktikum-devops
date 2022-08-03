package agent

import (
	"context"

	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

// These are the interfaces required by the Agent to work with metrics

// Worker bundles together common functionality for collectors/exporters
// Worker has a name and can be reserved for work or released when done
type Worker interface {
	// Name returns user-friendly name for the worker
	Name() string
	// MaxThreads returns max amount of parallel threads doing work at the same time
	MaxThreads() int
	// Reserve tries to reserve the Worker until it is available or the context is canceled
	Reserve(ctx context.Context) bool
	// Release is the opposite of Reserve(), it tries to make the Worker available instead of busy
	Release(ctx context.Context) bool
}

// MetricsCollector is a Worker that can collect metrics from somewhere
type MetricsCollector interface {
	Worker
	Collect(context.Context) ([]domain.Metric, error)
}

// MetricsExporter is a Worker that can export metrics to somewhere
type MetricsExporter interface {
	Worker
	Export(context.Context, []domain.Metric) error
}

// MetricsBufferer can buffer metrics in temporary storage before exporting
type MetricsBufferer interface {
	Buffer([]domain.Metric)
	Retrieve() []domain.Metric
	Flush()
}
