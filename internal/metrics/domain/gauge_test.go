package domain

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGaugeStringValue(t *testing.T) {
	tests := []struct {
		name   string
		metric *gauge
		want   string
	}{
		{
			name:   "generic test #1",
			metric: NewGauge(Alloc, 10.20),
			want:   "10.2",
		},
		{
			name:   "generic test #2",
			metric: NewGauge(Alloc, 123.456789),
			want:   "123.456789",
		},
		{
			name:   "zero",
			metric: NewGauge(Alloc, 0),
			want:   "0.0",
		},
		{
			name:   "negative",
			metric: NewGauge(Alloc, -100.5),
			want:   "-100.5",
		},
		{
			name:   "very big number",
			metric: NewGauge(Alloc, math.MaxFloat64),
			want: "17976931348623157081452742373170435679807056752584499659891747" +
				"6803157260780028538760589558632766878171540458953514382464234321326" +
				"8894641827684675467035375169860499105765512820762454900903893289440" +
				"7586850845513394230458323690322294816580855933212334827479782620414" +
				"4723168738177180919299881250404026184124858368.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.metric.StringValue())
		})
	}
}

func TestGaugeSet(t *testing.T) {
	tests := []struct {
		name  string
		have  *gauge
		value Gauge
	}{
		{
			name:  "generic test",
			have:  NewGauge(Alloc, 10.20),
			value: Gauge(5.5),
		},
		{
			name:  "generic test #2",
			have:  NewGauge(Alloc, 123.456789),
			value: Gauge(123.456789),
		},
		{
			name:  "zero",
			have:  NewGauge(Alloc, 0),
			value: Gauge(0.0),
		},
		{
			name:  "negative",
			have:  NewGauge(Alloc, -100.5),
			value: Gauge(-100.5),
		},
		{
			name:  "very big number",
			have:  NewGauge(Alloc, math.MaxFloat64),
			value: Gauge(math.MaxFloat64),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := tt.value
			err := tt.have.Set(tt.value)
			require.NoError(t, err)
			assert.Equal(t, want, tt.have.value)
		})
	}
}
