package collectors

import (
	"context"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics"
)

type Collector interface {
	GetName() string
	Collect(context.Context) ([]metrics.Metric, error)
	Ready() <-chan bool
}

type AbstractCollector struct {
	name  string
	ready chan bool
}

func (col *AbstractCollector) GetName() string {
	return col.name
}

func (col *AbstractCollector) Ready() <-chan bool {
	return col.ready
}

func (col *AbstractCollector) readyUp() {
	go func() {
		col.ready <- true
	}()
}
