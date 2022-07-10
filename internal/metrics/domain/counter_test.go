package domain

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounterStringValue(t *testing.T) {
	tests := []struct {
		name   string
		metric *counter
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
			assert.Equal(t, tt.want, tt.metric.StringValue())
		})
	}
}

func TestCounterUpdate(t *testing.T) {
	tests := []struct {
		name  string
		have  *counter
		value Counter
		want  Counter
	}{
		{
			name:  "generic test #1",
			have:  NewCounter(PollCount, 10),
			value: Counter(5),
			want:  Counter(15),
		},
		{
			name:  "generic test #2",
			have:  NewCounter(PollCount, 497285126312),
			value: Counter(12732172343243),
			want:  Counter(13229457469555),
		},
		{
			name:  "add zero",
			have:  NewCounter(PollCount, 123),
			value: Counter(0),
			want:  Counter(123),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.have.Update(tt.value)
			assert.Equal(t, tt.want, tt.have.value)
		})
	}
}
