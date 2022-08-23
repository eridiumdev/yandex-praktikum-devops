package http

import (
	"context"

	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

// These are the interfaces required for handling metrics requests

// MetricsRenderer should apply metrics to some template, resulting in renderable output
type MetricsRenderer interface {
	RenderList(list []domain.Metric) ([]byte, error)
}

// MetricsService should be able to perform common operations on metrics, such as updating and retrieving
type MetricsService interface {
	Update(ctx context.Context, metric domain.Metric) (domain.Metric, error)
	UpdateMany(ctx context.Context, metrics []domain.Metric) ([]domain.Metric, error)
	Get(ctx context.Context, name string) (m domain.Metric, found bool, err error)
	List(ctx context.Context) ([]domain.Metric, error)
}

// MetricsRequestResponseFactory can build various requests/responses for usage in the handler
type MetricsRequestResponseFactory interface {
	BuildUpdateMetricRequest(ctx context.Context, metric domain.Metric) domain.UpdateMetricRequest
	BuildUpdateMetricResponse(ctx context.Context, metric domain.Metric) domain.UpdateMetricResponse
	BuildUpdateBatchMetricRequest(ctx context.Context, metrics []domain.Metric) []domain.UpdateMetricRequest
	BuildUpdateBatchMetricResponse(ctx context.Context, metrics []domain.Metric) []domain.UpdateMetricResponse
	BuildGetMetricResponse(ctx context.Context, metric domain.Metric) domain.GetMetricResponse
}

// MetricsHasher can calculate hashes based on metric, and also check if provided hash matches calculated
type MetricsHasher interface {
	Hash(ctx context.Context, metric domain.Metric) string
	Check(ctx context.Context, metric domain.Metric, hash string) bool
}
