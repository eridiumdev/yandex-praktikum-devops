package metrics

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestCounterGetStringValue(t *testing.T) {
	tests := []struct {
		name   string
		metric CounterMetric
		want   string
	}{
		{
			name: "generic test #1",
			metric: CounterMetric{
				AbstractMetric: AbstractMetric{
					Name: PollCount,
				},
				Value: 10,
			},
			want: "10",
		},
		{
			name: "generic test #2",
			metric: CounterMetric{
				AbstractMetric: AbstractMetric{
					Name: PollCount,
				},
				Value: 100500,
			},
			want: "100500",
		},
		{
			name: "zero",
			metric: CounterMetric{
				AbstractMetric: AbstractMetric{
					Name: PollCount,
				},
				Value: 0,
			},
			want: "0",
		},
		{
			name: "very big number",
			metric: CounterMetric{
				AbstractMetric: AbstractMetric{
					Name: PollCount,
				},
				Value: math.MaxInt,
			},
			want: "9223372036854775807",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.metric.GetStringValue())
		})
	}
}
