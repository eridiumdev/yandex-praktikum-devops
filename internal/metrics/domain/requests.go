package domain

type genericMetric struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type UpdateMetricRequest struct {
	genericMetric
}

type UpdateMetricResponse struct {
	genericMetric
}

type GetMetricRequest struct {
	genericMetric
}

type GetMetricResponse struct {
	genericMetric
}

func PrepareUpdateMetricRequest(metric Metric) UpdateMetricRequest {
	return UpdateMetricRequest{
		populateGenericMetric(metric),
	}
}

func PrepareUpdateMetricResponse(metric Metric) UpdateMetricResponse {
	return UpdateMetricResponse{
		populateGenericMetric(metric),
	}
}

func PrepareGetMetricResponse(metric Metric) GetMetricResponse {
	return GetMetricResponse{
		populateGenericMetric(metric),
	}
}

func populateGenericMetric(metric Metric) genericMetric {
	result := genericMetric{
		ID:    metric.Name,
		MType: metric.Type,
	}
	switch metric.Type {
	case TypeCounter:
		val := int64(metric.Counter)
		result.Delta = &val
	case TypeGauge:
		val := float64(metric.Gauge)
		result.Value = &val
	}
	return result
}
