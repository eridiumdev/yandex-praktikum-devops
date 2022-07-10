package domain

import (
	"strconv"
	"strings"
)

type Gauge float64

type gauge struct {
	*abstractMetric
	value Gauge
}

func NewGauge(name string, value Gauge) *gauge {
	return &gauge{
		abstractMetric: &abstractMetric{
			name: name,
		},
		value: value,
	}
}

func (m *gauge) Type() string {
	return TypeGauge
}

func (m *gauge) Value() MetricValue {
	return m.value
}

func (m *gauge) StringValue() string {
	trimmed := strings.TrimRight(strconv.FormatFloat(float64(m.value), 'f', 6, 64), "0")
	if trimmed[len(trimmed) - 1] == '.' {
		// Add zero after decimal point
		// e.g. '10.000000' after trimming will be '10.' -> add '0' to become '10.0'
		return trimmed + "0"
	} else {
		return trimmed
	}
}

// Update resets current gauge value
func (m *gauge) Update(value MetricValue) {
	m.value = value.(Gauge)
}

// Copy creates a copy of gauge with same name/value
func (m *gauge) Copy() Metric {
	return NewGauge(m.name, m.value)
}
