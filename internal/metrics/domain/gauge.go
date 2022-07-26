package domain

import (
	"strconv"
	"strings"
)

type Gauge float64

func NewGauge(name string, value Gauge) Metric {
	return Metric{
		Name:  name,
		Type:  TypeGauge,
		Gauge: value,
	}
}

func (g Gauge) String() string {
	trimmed := strings.TrimRight(strconv.FormatFloat(float64(g), 'f', 6, 64), "0")
	if trimmed[len(trimmed)-1] == '.' {
		// Add zero after decimal point
		// e.g. '10.000000' after trimming will be '10.' -> add '0' to become '10.0'
		return trimmed + "0"
	}
	return trimmed
}

func (m Metric) IsGauge() bool {
	return m.Type == TypeGauge
}
