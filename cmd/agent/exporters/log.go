package exporters

import (
	"context"

	"eridiumdev/yandex-praktikum-go-devops/internal/logger"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics"
)

type LogExporter struct {
	*AbstractExporter
}

func NewLogExporter(name string) *LogExporter {
	exp := &LogExporter{
		AbstractExporter: &AbstractExporter{
			name:  name,
			ready: make(chan bool),
		},
	}
	exp.readyUp()
	return exp
}

func (exp *LogExporter) Export(ctx context.Context, mtx []metrics.Metric) error {
	defer func() {
		exp.readyUp()
	}()

	for _, metric := range mtx {
		logger.Infof("%s:%s (%s)", metric.GetName(), metric.GetStringValue(), metric.GetType())
	}
	return nil
}
