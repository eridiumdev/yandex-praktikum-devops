package metrics

import "strconv"

type Gauge float64

type GaugeMetric struct {
	AbstractMetric
	Value Gauge
}

func (m GaugeMetric) GetType() string {
	return TypeGauge
}

func (m GaugeMetric) GetValue() interface{} {
	return m.Value
}

func (m GaugeMetric) GetStringValue() string {
	return strconv.FormatFloat(float64(m.Value), 'f', 6, 64)
}
