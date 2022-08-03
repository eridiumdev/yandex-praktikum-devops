package agent

import (
	"context"
	"time"

	"eridiumdev/yandex-praktikum-go-devops/config"
	"eridiumdev/yandex-praktikum-go-devops/internal/common/logger"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

type Agent struct {
	collectInterval time.Duration
	exportInterval  time.Duration

	collectors []MetricsCollector
	exporters  []MetricsExporter
	bufferer   MetricsBufferer
}

func NewAgent(cfg *config.AgentConfig, bufferer MetricsBufferer) *Agent {
	return &Agent{
		collectInterval: cfg.CollectInterval,
		exportInterval:  cfg.ExportInterval,
		collectors:      []MetricsCollector{},
		exporters:       []MetricsExporter{},
		bufferer:        bufferer,
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
	ticker := time.NewTicker(a.collectInterval)
	for {
		select {
		case <-ticker.C:
			collectCycles++
			logger.New(ctx).Debugf("[agent] collecting cycle %d", collectCycles)
			for _, col := range a.collectors {
				go a.collectMetrics(ctx, col)
			}
		case <-ctx.Done():
			logger.New(ctx).Debugf("[agent] context cancelled, collecting stopped")
			return
		}
	}
}

func (a *Agent) StartExporting(ctx context.Context) {
	exportCycles := 0
	ticker := time.NewTicker(a.exportInterval)
	for {
		select {
		case <-ticker.C:
			exportCycles++
			logger.New(ctx).Debugf("[agent] exporting cycle %d", exportCycles)

			// Get current bufferer snapshot
			bufferSnapshot := a.bufferer.Retrieve()
			// Send metrics to exporters
			for _, exp := range a.exporters {
				go a.exportMetrics(ctx, exp, bufferSnapshot)
			}
			// Flush the bufferer after exporting
			a.bufferer.Flush()
		case <-ctx.Done():
			logger.New(ctx).Debugf("[agent] context cancelled, exporting stopped")
			return
		}
	}
}

func (a *Agent) Stop(ctx context.Context) {
	// Wait for collectors to finish their job
	// (= try to reserve all available worker threads)
	for _, col := range a.collectors {
		for i := 0; i < col.MaxThreads(); i++ {
			col.Reserve(ctx)
		}
	}
	// Wait for exporters to finish their job
	// (= try to reserve all available worker threads)
	for _, exp := range a.exporters {
		for i := 0; i < exp.MaxThreads(); i++ {
			exp.Reserve(ctx)
		}
	}
}

func (a *Agent) collectMetrics(ctx context.Context, col MetricsCollector) {
	reserveCtx, cancel := context.WithTimeout(ctx, a.collectInterval)
	defer cancel()

	// Try to reserve the collector
	ok := col.Reserve(reserveCtx)
	if !ok {
		logger.New(ctx).Errorf("[%s collector] timeout when collecting metrics: collector still busy", col.Name())
	}
	defer col.Release(ctx)

	logger.New(ctx).Debugf("[%s collector] start collecting metrics", col.Name())
	snapshot, err := col.Collect(ctx)
	if err != nil {
		logger.New(ctx).Errorf("[%s collector] error when collecting metrics: %s", col.Name(), err.Error())
	}
	a.bufferer.Buffer(snapshot)
	logger.New(ctx).Debugf("[%s collector] finish collecting metrics", col.Name())
}

func (a *Agent) exportMetrics(ctx context.Context, exp MetricsExporter, metrics []domain.Metric) {
	reserveCtx, cancel := context.WithTimeout(ctx, a.exportInterval)
	defer cancel()

	// Try to reserve the exporter
	ok := exp.Reserve(reserveCtx)
	if !ok {
		logger.New(ctx).Errorf("[%s exporter] timeout when exporting metrics: exporter still busy", exp.Name())
	}
	defer exp.Release(ctx)

	logger.New(ctx).Debugf("[%s exporter] start exporting metrics", exp.Name())
	err := exp.Export(ctx, metrics)
	if err != nil {
		logger.New(ctx).Errorf("[%s exporter] error when exporting metrics: %s", exp.Name(), err.Error())
	}
	logger.New(ctx).Debugf("[%s exporter] finish exporting metrics", exp.Name())
}
