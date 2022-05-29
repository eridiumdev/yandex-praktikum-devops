package main

import (
	"eridiumdev/yandex-praktikum-go-devops/cmd/agent/collectors"
	"eridiumdev/yandex-praktikum-go-devops/cmd/agent/exporters"
	"eridiumdev/yandex-praktikum-go-devops/internal/logger"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics"
	"time"
)

type Agent struct {
	collectInterval time.Duration
	exportInterval  time.Duration

	collectors []collectors.Collector
	exporters  []exporters.Exporter

	metricsChannel chan metrics.Metric
	metricsBuffer  []metrics.Metric
}

func NewAgent(collectInterval time.Duration, exportInterval time.Duration) *Agent {
	return &Agent{
		collectInterval: collectInterval,
		exportInterval:  exportInterval,
		collectors:      []collectors.Collector{},
		exporters:       []exporters.Exporter{},
		metricsChannel:  make(chan metrics.Metric),
		metricsBuffer:   []metrics.Metric{},
	}
}

func (a *Agent) AddCollector(col collectors.Collector) {
	a.collectors = append(a.collectors, col)
}

func (a *Agent) AddExporter(exp exporters.Exporter) {
	a.exporters = append(a.exporters, exp)
}

func (a *Agent) StartCollecting() {
	collectCycles := 0
	for range time.Tick(a.collectInterval) {
		collectCycles++
		logger.Debugf("[agent] collecting cycle %d", collectCycles)
		for _, col := range a.collectors {
			go a.collectMetrics(col)
		}
	}
}

func (a *Agent) StartExporting() {
	exportCycles := 0
	for range time.Tick(a.exportInterval) {
		exportCycles++
		logger.Debugf("[agent] exporting cycle %d", exportCycles)
		for _, exp := range a.exporters {
			go a.exportMetrics(exp, a.metricsBuffer)
		}
	}
}

func (a *Agent) StartBuffering() {
	for metric := range a.metricsChannel {
		found := false
		for i, m := range a.metricsBuffer {
			if m.GetName() == metric.GetName() {
				// Overwrite previous metric in buffer, if set
				a.metricsBuffer[i] = metric
				found = true
				break
			}
		}
		if !found {
			// Add metric to buffer for the first time
			a.metricsBuffer = append(a.metricsBuffer, metric)
		}
	}
}

func (a *Agent) collectMetrics(col collectors.Collector) {
	select {
	case <-col.Ready():
		logger.Debugf("[%s collector] start collecting metrics", col.GetName())
		snapshot, err := col.Collect()
		if err != nil {
			logger.Errorf("[%s collector] error when collecting metrics: %s", col.GetName(), err.Error())
		}
		for _, metric := range snapshot {
			a.metricsChannel <- metric
		}
		logger.Debugf("[%s collector] finish collecting metrics", col.GetName())
	case <-time.After(a.collectInterval):
		logger.Errorf("[%s collector] timeout when collecting metrics: collector not yet ready", col.GetName())
		return
	}
}

func (a *Agent) exportMetrics(exp exporters.Exporter, metrics []metrics.Metric) {
	select {
	case <-exp.Ready():
		logger.Debugf("[%s exporter] start exporting metrics", exp.GetName())
		err := exp.Export(metrics)
		if err != nil {
			logger.Errorf("[%s exporter] error when exporting metrics: %s", exp.GetName(), err.Error())
		}
		logger.Debugf("[%s exporter] finish exporting metrics", exp.GetName())
	case <-time.After(a.exportInterval):
		logger.Errorf("[%s exporter] timeout when exporting metrics: exporter not yet ready", exp.GetName())
		return
	}
}
