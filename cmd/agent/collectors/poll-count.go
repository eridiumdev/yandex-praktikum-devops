package collectors

import (
	"context"

	"eridiumdev/yandex-praktikum-go-devops/internal/metrics"
)

type PollCountCollector struct {
	*AbstractCollector
	pollCount metrics.Counter
}

func NewPollCountCollector(name string) *PollCountCollector {
	col := &PollCountCollector{
		AbstractCollector: &AbstractCollector{
			name:  name,
			ready: make(chan bool),
		},
	}
	col.readyUp()
	return col
}

func (col *PollCountCollector) Collect(ctx context.Context) ([]metrics.Metric, error) {
	defer func() {
		col.readyUp()
	}()

	return []metrics.Metric{
		metrics.CounterMetric{
			AbstractMetric: metrics.AbstractMetric{
				Name: metrics.PollCount,
			},
			Value: 1,
		},
	}, nil
}
