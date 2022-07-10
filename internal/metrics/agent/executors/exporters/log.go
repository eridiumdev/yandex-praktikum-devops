package exporters

import (
	"context"

	"eridiumdev/yandex-praktikum-go-devops/internal/commons/logger"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/agent/executors"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

type LogExporter struct {
	*executors.Executor
}

func NewLogExporter(name string) *LogExporter {
	exp := &LogExporter{
		Executor: executors.New(name),
	}
	exp.ReadyUp()
	return exp
}

func (exp *LogExporter) Export(ctx context.Context, mtx []domain.Metric) error {
	defer func() {
		exp.ReadyUp()
	}()

	for _, metric := range mtx {
		logger.Infof("%s:%s (%s)", metric.Name(), metric.StringValue(), metric.Type())
	}
	return nil
}
