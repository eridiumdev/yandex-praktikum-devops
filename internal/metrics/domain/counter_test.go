package domain

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounterString(t *testing.T) {
	tests := []struct {
		name   string
		metric Metric
		want   string
	}{
		{
			name:   "generic test #1",
			metric: NewCounter(PollCount, 10),
			want:   "10",
		},
		{
			name:   "generic test #2",
			metric: NewCounter(PollCount, 100500),
			want:   "100500",
		},
		{
			name:   "zero",
			metric: NewCounter(PollCount, 0),
			want:   "0",
		},
		{
			name:   "very big number",
			metric: NewCounter(PollCount, math.MaxInt),
			want:   "9223372036854775807",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.metric.Counter.String())
		})
	}
}
