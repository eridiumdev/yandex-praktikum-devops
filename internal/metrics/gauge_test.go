package metrics

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGaugeGetStringValue(t *testing.T) {
	tests := []struct {
		name   string
		metric GaugeMetric
		want   string
	}{
		{
			name: "generic test #1",
			metric: GaugeMetric{
				AbstractMetric: AbstractMetric{
					Name: Alloc,
				},
				Value: 10.20,
			},
			want: "10.2",
		},
		{
			name: "generic test #2",
			metric: GaugeMetric{
				AbstractMetric: AbstractMetric{
					Name: Alloc,
				},
				Value: 123.456789,
			},
			want: "123.456789",
		},
		{
			name: "zero",
			metric: GaugeMetric{
				AbstractMetric: AbstractMetric{
					Name: Alloc,
				},
				Value: 0,
			},
			want: "0.0",
		},
		{
			name: "negative",
			metric: GaugeMetric{
				AbstractMetric: AbstractMetric{
					Name: Alloc,
				},
				Value: -100.5,
			},
			want: "-100.5",
		},
		{
			name: "very big number",
			metric: GaugeMetric{
				AbstractMetric: AbstractMetric{
					Name: Alloc,
				},
				Value: math.MaxFloat64,
			},
			want: "179769313486231570814527423731704356798070567525844996598917476803157260780028538760589558632766878171540458953514382464234321326889464182768467546703537516986049910576551282076245490090389328944075868508455133942304583236903222948165808559332123348274797826204144723168738177180919299881250404026184124858368.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.metric.GetStringValue())
		})
	}
}
