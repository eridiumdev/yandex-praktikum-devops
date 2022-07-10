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
	Update(value MetricValue)
	Copy() Metric
}

type abstractMetric struct {
	name string
}

func (m *abstractMetric) Name() string {
	return m.name
}
