package main

import (
	"context"
	"time"

	"eridiumdev/yandex-praktikum-go-devops/internal/common/logger"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

type AgentSettings struct {
	CollectInterval time.Duration
	ExportInterval  time.Duration
}

type Agent struct {
	AgentSettings
	collectors []MetricsCollector
	exporters  []MetricsExporter
	bufferer   MetricsBufferer
}

func NewAgent(settings AgentSettings, bufferer MetricsBufferer) *Agent {
	return &Agent{
		AgentSettings: settings,
		collectors:    []MetricsCollector{},
		exporters:     []MetricsExporter{},
		bufferer:      bufferer,
	}
}

func (a *Agent) AddCollector(col MetricsCollector) {
	a.collectors = append(a.collectors, col)
}

func (a *Agent) AddExporter(exp MetricsExporter) {
	a.exporters = append(a.exporters, exp)
}

func (a *Agent) StartCollecting(ctx context.Context) {
	collectCycles := 0
	for {
		select {
		case <-time.Tick(a.CollectInterval):
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
		case <-time.Tick(a.ExportInterval):
			exportCycles++
			logger.Debugf("[agent] exporting cycle %d", exportCycles)

			// Get current bufferer snapshot
			bufferSnapshot := a.bufferer.Retrieve()
			// Send metrics to exporters
			for _, exp := range a.exporters {
				go a.exportMetrics(ctx, exp, bufferSnapshot)
			}
			// Flush the bufferer after exporting
			a.bufferer.Flush()
		case <-ctx.Done():
			logger.Debugf("[agent] context cancelled, exporting stopped")
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

func (a *Agent) collectMetrics(ctx context.Context, col MetricsCollector) {
	select {
	case <-col.Ready():
		logger.Debugf("[%s collector] start collecting metrics", col.Name())
		snapshot, err := col.Collect(ctx)
		if err != nil {
			logger.Errorf("[%s collector] error when collecting metrics: %s", col.Name(), err.Error())
		}
		a.bufferer.Buffer(snapshot)
		logger.Debugf("[%s collector] finish collecting metrics", col.Name())
	case <-time.After(a.CollectInterval):
		logger.Errorf("[%s collector] timeout when collecting metrics: collector not yet ready", col.Name())
	case <-ctx.Done():
		logger.Debugf("[%s collector] context cancelled, skip collecting", col.Name())
	}
}

func (a *Agent) exportMetrics(ctx context.Context, exp MetricsExporter, metrics []domain.Metric) {
	select {
	case <-exp.Ready():
		logger.Debugf("[%s exporter] start exporting metrics", exp.Name())
		err := exp.Export(ctx, metrics)
		if err != nil {
			logger.Errorf("[%s exporter] error when exporting metrics: %s", exp.Name(), err.Error())
		}
		logger.Debugf("[%s exporter] finish exporting metrics", exp.Name())
	case <-time.After(a.ExportInterval):
		logger.Errorf("[%s exporter] timeout when exporting metrics: exporter not yet ready", exp.Name())
	case <-ctx.Done():
		logger.Debugf("[%s exporter] context cancelled, skip exporting", exp.Name())
	}
}
