package metrics

import (
	"strconv"
	"strings"
)

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
	trimmed := strings.TrimRight(strconv.FormatFloat(float64(m.Value), 'f', 6, 64), "0")
	if trimmed[len(trimmed) - 1] == '.' {
		return trimmed + "0"
	} else {
		return trimmed
	}
}
