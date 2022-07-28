package domain

import (
	"errors"
	"strconv"
)

type Counter int64

type counter struct {
	*abstractMetric
	value Counter
}

var ErrMetricValueNotCounter = errors.New("metric value is not a Counter")

func NewCounter(name string, value Counter) *counter {
	return &counter{
		abstractMetric: &abstractMetric{
			name: name,
		},
		value: value,
	}
}

func (m *counter) Type() string {
	return TypeCounter
}

func (m *counter) Value() MetricValue {
	return m.value
}

func (m *counter) StringValue() string {
	return strconv.FormatInt(int64(m.value), 10)
}

func (m *counter) Add(value MetricValue) error {
	val, ok := value.(Counter)
	if !ok {
		return ErrMetricValueNotCounter
	}
	m.value += val
	return nil
}

func (m *counter) Set(value MetricValue) error {
	val, ok := value.(Counter)
	if !ok {
		return ErrMetricValueNotCounter
	}
	m.value = val
	return nil
}

// Copy creates a copy of counter with same name/value
func (m *counter) Copy() Metric {
	return NewCounter(m.name, m.value)
}
