package domain

const (
	TypeGauge   = "gauge"
	TypeCounter = "counter"
)

type Metric struct {
	Name    string
	Type    string
	Counter Counter
	Gauge   Gauge
}

func (m Metric) StringValue() string {
	switch m.Type {
	case TypeCounter:
		return m.Counter.String()
	case TypeGauge:
		return m.Gauge.String()
	default:
		return ""
	}
}

func IsValidMetricType(metricType string) bool {
	for _, possible := range []string{TypeCounter, TypeGauge} {
		if metricType == possible {
			return true
		}
	}
	return false
}
