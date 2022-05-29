package metrics

import "strconv"

type Counter int64

type CounterMetric struct {
	AbstractMetric
	Value Counter
}

func (m CounterMetric) GetType() string {
	return TypeCounter
}

func (m CounterMetric) GetValue() interface{} {
	return m.Value
}

func (m CounterMetric) GetStringValue() string {
	return strconv.FormatInt(int64(m.Value), 10)
}
