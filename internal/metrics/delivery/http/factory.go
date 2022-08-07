package http

import (
	"context"

	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

type requestResponseFactory struct {
	hasher MetricsHasher
}

func NewRequestResponseFactory(hasher MetricsHasher) *requestResponseFactory {
	return &requestResponseFactory{
		hasher: hasher,
	}
}

func (f *requestResponseFactory) BuildUpdateMetricRequest(
	ctx context.Context,
	metric domain.Metric,
) domain.UpdateMetricRequest {
	return domain.UpdateMetricRequest{
		GenericMetric: f.populateGenericMetric(ctx, metric),
	}
}

func (f *requestResponseFactory) BuildUpdateMetricResponse(
	ctx context.Context,
	metric domain.Metric,
) domain.UpdateMetricResponse {
	return domain.UpdateMetricResponse{
		GenericMetric: f.populateGenericMetric(ctx, metric),
	}
}

func (f *requestResponseFactory) BuildGetMetricResponse(
	ctx context.Context,
	metric domain.Metric,
) domain.GetMetricResponse {
	return domain.GetMetricResponse{
		GenericMetric: f.populateGenericMetric(ctx, metric),
	}
}

func (f *requestResponseFactory) populateGenericMetric(ctx context.Context, metric domain.Metric) domain.GenericMetric {
	result := domain.GenericMetric{
		ID:    metric.Name,
		MType: metric.Type,
	}
	switch metric.Type {
	case domain.TypeCounter:
		val := int64(metric.Counter)
		result.Delta = &val
	case domain.TypeGauge:
		val := float64(metric.Gauge)
		result.Value = &val
	}
	result.Hash = f.hasher.Hash(ctx, metric)
	return result
}
