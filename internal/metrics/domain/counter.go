package domain

import "strconv"

type Counter int64

type counter struct {
	*abstractMetric
	value Counter
}

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

// Update increments counter value
func (m *counter) Update(value MetricValue) {
	m.value += value.(Counter)
}

// Copy creates a copy of counter with same name/value
func (m *counter) Copy() Metric {
	return NewCounter(m.name, m.value)
}
