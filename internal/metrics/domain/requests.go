package domain

type GenericMetric struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	Hash  string   `json:"hash,omitempty"`
}

type UpdateMetricRequest struct {
	GenericMetric
}

type UpdateMetricResponse struct {
	GenericMetric
}

type GetMetricRequest struct {
	GenericMetric
}

type GetMetricResponse struct {
	GenericMetric
}

func (g GenericMetric) TranslateToMetric() Metric {
	metric := Metric{
		Name: g.ID,
		Type: g.MType,
	}
	if g.Delta != nil {
		metric.Counter = Counter(*g.Delta)
	}
	if g.Value != nil {
		metric.Gauge = Gauge(*g.Value)
	}
	return metric
}
