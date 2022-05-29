package collectors

import "eridiumdev/yandex-praktikum-go-devops/internal/metrics"

type Collector interface {
	GetName() string
	Collect() ([]metrics.Metric, error)
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

func (col *AbstractCollector) makeReady() {
	go func() {
		col.ready <- true
	}()
}
