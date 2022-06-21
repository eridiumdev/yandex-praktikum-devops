package exporters

import (
	"context"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics"
)

type Exporter interface {
	GetName() string
	Export(context.Context, []metrics.Metric) error
	Ready() <-chan bool
}

type AbstractExporter struct {
	name  string
	ready chan bool
}

func (exp *AbstractExporter) GetName() string {
	return exp.name
}

func (exp *AbstractExporter) Ready() <-chan bool {
	return exp.ready
}

func (exp *AbstractExporter) readyUp() {
	go func() {
		exp.ready <- true
	}()
}
