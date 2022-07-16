package domain

import (
	"errors"
	"strconv"
	"strings"
)

type Gauge float64

type gauge struct {
	*abstractMetric
	value Gauge
}

var ErrMetricValueNotGauge = errors.New("metric value is not a Counter")

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
	if trimmed[len(trimmed)-1] == '.' {
		// Add zero after decimal point
		// e.g. '10.000000' after trimming will be '10.' -> add '0' to become '10.0'
		return trimmed + "0"
	}
	return trimmed
}

func (m *gauge) Add(value MetricValue) error {
	val, ok := value.(Gauge)
	if !ok {
		return ErrMetricValueNotGauge
	}
	m.value += val
	return nil
}

func (m *gauge) Set(value MetricValue) error {
	val, ok := value.(Gauge)
	if !ok {
		return ErrMetricValueNotGauge
	}
	m.value = val
	return nil
}

// Copy creates a copy of gauge with same name/value
func (m *gauge) Copy() Metric {
	return NewGauge(m.name, m.value)
}
