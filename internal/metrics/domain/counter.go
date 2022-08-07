package domain

import "strconv"

type Counter int64

func NewCounter(name string, value Counter) Metric {
	return Metric{
		Name:    name,
		Type:    TypeCounter,
		Counter: value,
	}
}

func (c Counter) String() string {
	return strconv.FormatInt(int64(c), 10)
}

func (m Metric) IsCounter() bool {
	return m.Type == TypeCounter
}
