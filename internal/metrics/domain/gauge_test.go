package domain

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGaugeStringValue(t *testing.T) {
	tests := []struct {
		name   string
		metric Metric
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
			assert.Equal(t, tt.want, tt.metric.Gauge.String())
		})
	}
}
