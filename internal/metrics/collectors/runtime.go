package collectors

import (
	"context"
	"runtime"

	"eridiumdev/yandex-praktikum-go-devops/internal/common/worker"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

type runtimeCollector struct {
	*worker.Worker
}

func NewRuntimeCollector(name string) *runtimeCollector {
	col := &runtimeCollector{
		Worker: worker.New(name, 1),
	}
	return col
}

func (col *runtimeCollector) Collect(ctx context.Context) ([]domain.Metric, error) {
	return col.getRuntimeSnapshot(), nil
}

func (col *runtimeCollector) getRuntimeSnapshot() []domain.Metric {
	stats := &runtime.MemStats{}
	runtime.ReadMemStats(stats)

	return []domain.Metric{
		domain.NewGauge(domain.Alloc, domain.Gauge(stats.Alloc)),
		domain.NewGauge(domain.BuckHashSys, domain.Gauge(stats.BuckHashSys)),
		domain.NewGauge(domain.Frees, domain.Gauge(stats.Frees)),
		domain.NewGauge(domain.GCCPUFraction, domain.Gauge(stats.GCCPUFraction)),
		domain.NewGauge(domain.GCSys, domain.Gauge(stats.GCSys)),
		domain.NewGauge(domain.HeapAlloc, domain.Gauge(stats.HeapAlloc)),
		domain.NewGauge(domain.HeapIdle, domain.Gauge(stats.HeapIdle)),
		domain.NewGauge(domain.HeapInuse, domain.Gauge(stats.HeapInuse)),
		domain.NewGauge(domain.HeapObjects, domain.Gauge(stats.HeapObjects)),
		domain.NewGauge(domain.HeapReleased, domain.Gauge(stats.HeapReleased)),
		domain.NewGauge(domain.HeapSys, domain.Gauge(stats.HeapSys)),
		domain.NewGauge(domain.LastGC, domain.Gauge(stats.LastGC)),
		domain.NewGauge(domain.Lookups, domain.Gauge(stats.Lookups)),
		domain.NewGauge(domain.MCacheInuse, domain.Gauge(stats.MCacheInuse)),
		domain.NewGauge(domain.MCacheSys, domain.Gauge(stats.MCacheSys)),
		domain.NewGauge(domain.MSpanInuse, domain.Gauge(stats.MSpanInuse)),
		domain.NewGauge(domain.MSpanSys, domain.Gauge(stats.MSpanSys)),
		domain.NewGauge(domain.Mallocs, domain.Gauge(stats.Mallocs)),
		domain.NewGauge(domain.NextGC, domain.Gauge(stats.NextGC)),
		domain.NewGauge(domain.NumForcedGC, domain.Gauge(stats.NumForcedGC)),
		domain.NewGauge(domain.NumGC, domain.Gauge(stats.NumGC)),
		domain.NewGauge(domain.OtherSys, domain.Gauge(stats.OtherSys)),
		domain.NewGauge(domain.PauseTotalNs, domain.Gauge(stats.PauseTotalNs)),
		domain.NewGauge(domain.StackInuse, domain.Gauge(stats.StackInuse)),
		domain.NewGauge(domain.StackSys, domain.Gauge(stats.StackSys)),
		domain.NewGauge(domain.Sys, domain.Gauge(stats.Sys)),
		domain.NewGauge(domain.TotalAlloc, domain.Gauge(stats.TotalAlloc)),
	}
}
