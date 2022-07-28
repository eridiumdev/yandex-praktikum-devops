package collectors

import (
	"context"

	"eridiumdev/yandex-praktikum-go-devops/internal/commons/executor"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

type pollCountCollector struct {
	*executor.Executor
}

func NewPollCountCollector(name string) *pollCountCollector {
	col := &pollCountCollector{
		Executor: executor.New(name),
	}
	col.ReadyUp()
	return col
}

func (col *pollCountCollector) Collect(ctx context.Context) ([]domain.Metric, error) {
	defer func() {
		col.ReadyUp()
	}()

	return []domain.Metric{
		domain.NewCounter(domain.PollCount, 1),
	}, nil
}
