package collectors

import (
	"context"

	"eridiumdev/yandex-praktikum-go-devops/internal/common/worker"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

type pollCountCollector struct {
	*worker.Worker
}

func NewPollCountCollector(name string) *pollCountCollector {
	col := &pollCountCollector{
		Worker: worker.New(name, 1),
	}
	return col
}

func (col *pollCountCollector) Collect(ctx context.Context) ([]domain.Metric, error) {
	return []domain.Metric{
		domain.NewCounter(domain.PollCount, 1),
	}, nil
}
