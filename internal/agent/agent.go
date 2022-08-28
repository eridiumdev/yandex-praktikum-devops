package agent

import (
	"context"
	"time"

	"eridiumdev/yandex-praktikum-go-devops/config"
	"eridiumdev/yandex-praktikum-go-devops/internal/common/logger"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

type agent struct {
	cfg        *config.AgentConfig
	bufferer   MetricsBufferer
	collectors []MetricsCollector
	exporters  []MetricsExporter
}

func newAgent(cfg *config.AgentConfig, bufferer MetricsBufferer) *agent {
	return &agent{
		cfg:        cfg,
		bufferer:   bufferer,
		collectors: []MetricsCollector{},
		exporters:  []MetricsExporter{},
	}
}

func (a *agent) LogSource() string {
	return "agent"
}

func (a *agent) AddCollector(col MetricsCollector) {
	a.collectors = append(a.collectors, col)
}

func (a *agent) AddExporter(exp MetricsExporter) {
	a.exporters = append(a.exporters, exp)
}

func (a *agent) StartCollecting(ctx context.Context) {
	collectCycles := 0
	ticker := time.NewTicker(a.cfg.CollectInterval)
	for {
		select {
		case <-ticker.C:
			collectCycles++
			logger.New(ctx).Src(a).Debugf("collecting cycle %d", collectCycles)
			for _, col := range a.collectors {
				go a.collectMetrics(ctx, col)
			}
		case <-ctx.Done():
			logger.New(ctx).Src(a).Debugf("context cancelled, collecting stopped")
			return
		}
	}
}

func (a *agent) StartExporting(ctx context.Context) {
	exportCycles := 0
	ticker := time.NewTicker(a.cfg.ExportInterval)
	for {
		select {
		case <-ticker.C:
			exportCycles++
			logger.New(ctx).Src(a).Debugf("exporting cycle %d", exportCycles)

			// Get current bufferer snapshot
			bufferSnapshot := a.bufferer.Retrieve()
			// Send metrics to exporters
			for _, exp := range a.exporters {
				go a.exportMetrics(ctx, exp, bufferSnapshot)
			}
			// Flush the bufferer after exporting
			a.bufferer.Flush()
		case <-ctx.Done():
			logger.New(ctx).Src(a).Debugf("context cancelled, exporting stopped")
			return
		}
	}
}

func (a *agent) collectMetrics(ctx context.Context, col MetricsCollector) {
	ctx = logger.Enrich(ctx, "collector", col.Name())
	reserveCtx, cancel := context.WithTimeout(ctx, a.cfg.CollectInterval)
	defer cancel()

	// Try to reserve the collector
	err := col.Reserve(reserveCtx)
	if err != nil {
		logger.New(ctx).Src(a).Err(err).Errorf("failed to reserve collector")
		return
	}
	defer func() {
		err := col.Release(ctx)
		if err != nil {
			logger.New(ctx).Src(a).Err(err).Errorf("failed to release collector")
		}
	}()

	logger.New(ctx).Src(a).Debugf("start collecting metrics")
	snapshot, err := col.Collect(ctx)
	if err != nil {
		logger.New(ctx).Src(a).Err(err).Errorf("error when collecting metrics")
	}
	a.bufferer.Buffer(snapshot)
	logger.New(ctx).Src(a).Debugf("finish collecting metrics")
}

func (a *agent) exportMetrics(ctx context.Context, exp MetricsExporter, metrics []domain.Metric) {
	reserveCtx, cancel := context.WithTimeout(ctx, a.cfg.ExportInterval)
	defer cancel()

	// Try to reserve the exporter
	err := exp.Reserve(reserveCtx)
	if err != nil {
		logger.New(ctx).Src(a).Field("exporter", exp.Name()).Err(err).Errorf("failed to reserve exporter")
		return
	}
	defer func() {
		err := exp.Release(ctx)
		if err != nil {
			logger.New(ctx).Src(a).Field("exporter", exp.Name()).Err(err).Errorf("failed to release exporter")
		}
	}()

	logger.New(ctx).Src(a).Field("exporter", exp.Name()).Debugf("start exporting metrics")
	err = exp.Export(ctx, metrics)
	if err != nil {
		logger.New(ctx).Src(a).Field("exporter", exp.Name()).Err(err).Errorf("error when exporting metrics")
	}
	logger.New(ctx).Src(a).Field("exporter", exp.Name()).Debugf("finish exporting metrics")
}
