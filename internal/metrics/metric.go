package metrics

const (
	TypeGauge   = "gauge"
	TypeCounter = "counter"
)

type Metric interface {
	GetName() string
	GetType() string
	GetValue() interface{}
	GetStringValue() string
}

type AbstractMetric struct {
	Name string
}

func (m AbstractMetric) GetName() string {
	return m.Name
}
