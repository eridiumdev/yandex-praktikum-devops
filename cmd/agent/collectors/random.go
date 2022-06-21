package collectors

import (
	"context"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics"
	"errors"
	"math/rand"
	"time"
)

type RandomCollector struct {
	*AbstractCollector
	generator      *rand.Rand
	randomValueMin int
	randomValueMax int
}

var (
	ErrNegativeNumber = errors.New("randomValueMin and randomValueMax cannot be negative")
	ErrMinOverMax     = errors.New("randomValueMin cannot be bigger than randomValueMax")
)

func NewRandomCollector(name string, randomValueMin, randomValueMax int) (*RandomCollector, error) {
	if randomValueMin < 0 || randomValueMax < 0 {
		return nil, ErrNegativeNumber
	}
	if randomValueMin > randomValueMax {
		return nil, ErrMinOverMax
	}

	col := &RandomCollector{
		AbstractCollector: &AbstractCollector{
			name:  name,
			ready: make(chan bool),
		},
		generator:      rand.New(rand.NewSource(time.Now().UnixNano())),
		randomValueMin: randomValueMin,
		randomValueMax: randomValueMax,
	}
	col.makeReady()
	return col, nil
}

func (col *RandomCollector) Collect(ctx context.Context) ([]metrics.Metric, error) {
	defer func() {
		col.ready <- true
	}()

	randomValue := col.generator.Intn((col.randomValueMax-col.randomValueMin)+1) + col.randomValueMin

	return []metrics.Metric{
		metrics.GaugeMetric{
			AbstractMetric: metrics.AbstractMetric{
				Name: metrics.RandomValue,
			},
			Value: metrics.Gauge(randomValue),
		},
	}, nil
}
