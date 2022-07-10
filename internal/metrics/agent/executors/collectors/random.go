package collectors

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/agent/executors"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

type randomCollector struct {
	*executors.Executor
	generator      *rand.Rand
	randomValueMin int
	randomValueMax int
}

var (
	ErrNegativeNumber = errors.New("randomValueMin and randomValueMax cannot be negative")
	ErrMinOverMax     = errors.New("randomValueMin cannot be bigger than randomValueMax")
)

func NewRandomCollector(name string, randomValueMin, randomValueMax int) (*randomCollector, error) {
	if randomValueMin < 0 || randomValueMax < 0 {
		return nil, ErrNegativeNumber
	}
	if randomValueMin > randomValueMax {
		return nil, ErrMinOverMax
	}

	col := &randomCollector{
		Executor:       executors.New(name),
		generator:      rand.New(rand.NewSource(time.Now().UnixNano())),
		randomValueMin: randomValueMin,
		randomValueMax: randomValueMax,
	}
	col.ReadyUp()
	return col, nil
}

func (col *randomCollector) Collect(ctx context.Context) ([]domain.Metric, error) {
	defer func() {
		col.ReadyUp()
	}()

	randomValue := col.generator.Intn((col.randomValueMax-col.randomValueMin)+1) + col.randomValueMin

	return []domain.Metric{
		domain.NewGauge(domain.RandomValue, domain.Gauge(randomValue)),
	}, nil
}
