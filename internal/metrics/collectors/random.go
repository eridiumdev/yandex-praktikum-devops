package collectors

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"eridiumdev/yandex-praktikum-go-devops/config"
	"eridiumdev/yandex-praktikum-go-devops/internal/common/worker"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

type randomCollector struct {
	*worker.Worker
	generator      *rand.Rand
	randomValueMin int
	randomValueMax int
}

var (
	ErrNegativeNumber = errors.New("randomValueMin and randomValueMax cannot be negative")
	ErrMinOverMax     = errors.New("randomValueMin cannot be bigger than randomValueMax")
)

func NewRandomCollector(name string, cfg config.RandomExporterConfig) (*randomCollector, error) {
	if cfg.Min < 0 || cfg.Max < 0 {
		return nil, ErrNegativeNumber
	}
	if cfg.Min > cfg.Max {
		return nil, ErrMinOverMax
	}

	col := &randomCollector{
		Worker:         worker.New(name, 1),
		generator:      rand.New(rand.NewSource(time.Now().UnixNano())),
		randomValueMin: cfg.Min,
		randomValueMax: cfg.Max,
	}
	return col, nil
}

func (col *randomCollector) Collect(ctx context.Context) ([]domain.Metric, error) {
	randomValue := col.generator.Intn((col.randomValueMax-col.randomValueMin)+1) + col.randomValueMin

	return []domain.Metric{
		domain.NewGauge(domain.RandomValue, domain.Gauge(randomValue)),
	}, nil
}
