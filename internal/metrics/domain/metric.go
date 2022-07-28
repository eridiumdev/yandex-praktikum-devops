package domain

const (
	TypeGauge   = "gauge"
	TypeCounter = "counter"
)

type MetricValue interface{}

type Metric interface {
	Name() string
	Type() string
	Value() MetricValue
	StringValue() string
	Add(value MetricValue) error
	Set(value MetricValue) error
	Copy() Metric
}

type abstractMetric struct {
	name string
}

func (m *abstractMetric) Name() string {
	return m.name
}

func IsValidMetricType(metricType string) bool {
	for _, possible := range []string{TypeCounter, TypeGauge} {
		if metricType == possible {
			return true
		}
	}
	return false
}
