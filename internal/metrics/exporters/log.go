package exporters

import (
	"context"

	"eridiumdev/yandex-praktikum-go-devops/internal/common/logger"
	"eridiumdev/yandex-praktikum-go-devops/internal/common/worker"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

type LogExporter struct {
	*worker.Worker
}

func NewLogExporter(name string) *LogExporter {
	exp := &LogExporter{
		Worker: worker.New(name, 1),
	}
	return exp
}

func (exp *LogExporter) Export(ctx context.Context, mtx []domain.Metric) error {
	for _, metric := range mtx {
		logger.New(ctx).Infof("%s:%s (%s)", metric.Name, metric.StringValue(), metric.Type)
	}
	return nil
}
