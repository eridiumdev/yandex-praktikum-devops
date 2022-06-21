package main

import (
	"context"
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

func (a *Agent) StartCollecting(ctx context.Context) {
	collectCycles := 0
	for {
		select {
		case <-time.Tick(a.collectInterval):
			collectCycles++
			logger.Debugf("[agent] collecting cycle %d", collectCycles)
			for _, col := range a.collectors {
				go a.collectMetrics(ctx, col)
			}
		case <-ctx.Done():
			logger.Debugf("[agent] context cancelled, collecting stopped")
			return
		}
	}
}

func (a *Agent) StartExporting(ctx context.Context) {
	exportCycles := 0
	for {
		select {
		case <-time.Tick(a.exportInterval):
			exportCycles++
			logger.Debugf("[agent] exporting cycle %d", exportCycles)
			for _, exp := range a.exporters {
				go a.exportMetrics(ctx, exp, a.metricsBuffer)
			}
		case <-ctx.Done():
			logger.Debugf("[agent] context cancelled, exporting stopped")
			return
		}
	}
}

func (a *Agent) StartBuffering(ctx context.Context) {
	for {
		select {
		case metric := <-a.metricsChannel:
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
		case <-ctx.Done():
			logger.Debugf("[agent] context cancelled, buffering stopped")
			return
		}
	}
}

func (a *Agent) Stop() {
	// Wait for collectors to finish their job
	for _, col := range a.collectors {
		<-col.Ready()
	}
	// Wait for exporters to finish their job
	for _, exp := range a.exporters {
		<-exp.Ready()
	}
}

func (a *Agent) collectMetrics(ctx context.Context, col collectors.Collector) {
	select {
	case <-col.Ready():
		logger.Debugf("[%s collector] start collecting metrics", col.GetName())
		snapshot, err := col.Collect(ctx)
		if err != nil {
			logger.Errorf("[%s collector] error when collecting metrics: %s", col.GetName(), err.Error())
		}
		for _, metric := range snapshot {
			a.metricsChannel <- metric
		}
		logger.Debugf("[%s collector] finish collecting metrics", col.GetName())
	case <-time.After(a.collectInterval):
		logger.Errorf("[%s collector] timeout when collecting metrics: collector not yet ready", col.GetName())
	case <-ctx.Done():
		logger.Debugf("[%s collector] context cancelled, skip collecting", col.GetName())
	}
}

func (a *Agent) exportMetrics(ctx context.Context, exp exporters.Exporter, metrics []metrics.Metric) {
	select {
	case <-exp.Ready():
		logger.Debugf("[%s exporter] start exporting metrics", exp.GetName())
		err := exp.Export(ctx, metrics)
		if err != nil {
			logger.Errorf("[%s exporter] error when exporting metrics: %s", exp.GetName(), err.Error())
		}
		logger.Debugf("[%s exporter] finish exporting metrics", exp.GetName())
	case <-time.After(a.exportInterval):
		logger.Errorf("[%s exporter] timeout when exporting metrics: exporter not yet ready", exp.GetName())
	case <-ctx.Done():
		logger.Debugf("[%s exporter] context cancelled, skip exporting", exp.GetName())
	}
}
