package collectors

import (
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics"
	"runtime"
)

type RuntimeCollector struct {
	*AbstractCollector
	pollCount metrics.Counter
}

func NewRuntimeCollector(name string) *RuntimeCollector {
	col := &RuntimeCollector{
		AbstractCollector: &AbstractCollector{
			name:  name,
			ready: make(chan bool),
		},
		pollCount: 0,
	}
	col.makeReady()
	return col
}

func (col *RuntimeCollector) Collect() ([]metrics.Metric, error) {
	defer func() {
		col.makeReady()
	}()
	// Get runtime metrics
	snapshot := col.getRuntimeSnapshot()

	// Increment pollCount after each poll
	col.pollCount++

	// Add pollCount to snapshot
	snapshot = append(snapshot, &metrics.CounterMetric{
		AbstractMetric: metrics.AbstractMetric{
			Name: metrics.PollCount,
		},
		Value: col.pollCount,
	})

	return snapshot, nil
}

func (col *RuntimeCollector) getRuntimeSnapshot() []metrics.Metric {
	stats := &runtime.MemStats{}
	runtime.ReadMemStats(stats)

	return []metrics.Metric{
		metrics.GaugeMetric{
			AbstractMetric: metrics.AbstractMetric{
				Name: metrics.Alloc,
			},
			Value: metrics.Gauge(stats.Alloc),
		},
		metrics.GaugeMetric{
			AbstractMetric: metrics.AbstractMetric{
				Name: metrics.BuckHashSys,
			},
			Value: metrics.Gauge(stats.BuckHashSys),
		},
		metrics.GaugeMetric{
			AbstractMetric: metrics.AbstractMetric{
				Name: metrics.Frees,
			},
			Value: metrics.Gauge(stats.Frees),
		},
		metrics.GaugeMetric{
			AbstractMetric: metrics.AbstractMetric{
				Name: metrics.GCCPUFraction,
			},
			Value: metrics.Gauge(stats.GCCPUFraction),
		},
		metrics.GaugeMetric{
			AbstractMetric: metrics.AbstractMetric{
				Name: metrics.GCSys,
			},
			Value: metrics.Gauge(stats.GCSys),
		},
		metrics.GaugeMetric{
			AbstractMetric: metrics.AbstractMetric{
				Name: metrics.HeapAlloc,
			},
			Value: metrics.Gauge(stats.HeapAlloc),
		},
		metrics.GaugeMetric{
			AbstractMetric: metrics.AbstractMetric{
				Name: metrics.HeapIdle,
			},
			Value: metrics.Gauge(stats.HeapIdle),
		},
		metrics.GaugeMetric{
			AbstractMetric: metrics.AbstractMetric{
				Name: metrics.HeapInuse,
			},
			Value: metrics.Gauge(stats.HeapInuse),
		},
		metrics.GaugeMetric{
			AbstractMetric: metrics.AbstractMetric{
				Name: metrics.HeapObjects,
			},
			Value: metrics.Gauge(stats.HeapObjects),
		},
		metrics.GaugeMetric{
			AbstractMetric: metrics.AbstractMetric{
				Name: metrics.HeapReleased,
			},
			Value: metrics.Gauge(stats.HeapReleased),
		},
		metrics.GaugeMetric{
			AbstractMetric: metrics.AbstractMetric{
				Name: metrics.HeapSys,
			},
			Value: metrics.Gauge(stats.HeapSys),
		},
		metrics.GaugeMetric{
			AbstractMetric: metrics.AbstractMetric{
				Name: metrics.LastGC,
			},
			Value: metrics.Gauge(stats.LastGC),
		},
		metrics.GaugeMetric{
			AbstractMetric: metrics.AbstractMetric{
				Name: metrics.Lookups,
			},
			Value: metrics.Gauge(stats.Lookups),
		},
		metrics.GaugeMetric{
			AbstractMetric: metrics.AbstractMetric{
				Name: metrics.MCacheInuse,
			},
			Value: metrics.Gauge(stats.MCacheInuse),
		},
		metrics.GaugeMetric{
			AbstractMetric: metrics.AbstractMetric{
				Name: metrics.MCacheSys,
			},
			Value: metrics.Gauge(stats.MCacheSys),
		},
		metrics.GaugeMetric{
			AbstractMetric: metrics.AbstractMetric{
				Name: metrics.MSpanInuse,
			},
			Value: metrics.Gauge(stats.MSpanInuse),
		},
		metrics.GaugeMetric{
			AbstractMetric: metrics.AbstractMetric{
				Name: metrics.MSpanSys,
			},
			Value: metrics.Gauge(stats.MSpanSys),
		},
		metrics.GaugeMetric{
			AbstractMetric: metrics.AbstractMetric{
				Name: metrics.Mallocs,
			},
			Value: metrics.Gauge(stats.Mallocs),
		},
		metrics.GaugeMetric{
			AbstractMetric: metrics.AbstractMetric{
				Name: metrics.NextGC,
			},
			Value: metrics.Gauge(stats.NextGC),
		},
		metrics.GaugeMetric{
			AbstractMetric: metrics.AbstractMetric{
				Name: metrics.NumForcedGC,
			},
			Value: metrics.Gauge(stats.NumForcedGC),
		},
		metrics.GaugeMetric{
			AbstractMetric: metrics.AbstractMetric{
				Name: metrics.NumGC,
			},
			Value: metrics.Gauge(stats.NumGC),
		},
		metrics.GaugeMetric{
			AbstractMetric: metrics.AbstractMetric{
				Name: metrics.OtherSys,
			},
			Value: metrics.Gauge(stats.OtherSys),
		},
		metrics.GaugeMetric{
			AbstractMetric: metrics.AbstractMetric{
				Name: metrics.PauseTotalNs,
			},
			Value: metrics.Gauge(stats.PauseTotalNs),
		},
		metrics.GaugeMetric{
			AbstractMetric: metrics.AbstractMetric{
				Name: metrics.StackInuse,
			},
			Value: metrics.Gauge(stats.StackInuse),
		},
		metrics.GaugeMetric{
			AbstractMetric: metrics.AbstractMetric{
				Name: metrics.StackSys,
			},
			Value: metrics.Gauge(stats.StackSys),
		},
		metrics.GaugeMetric{
			AbstractMetric: metrics.AbstractMetric{
				Name: metrics.Sys,
			},
			Value: metrics.Gauge(stats.Sys),
		},
		metrics.GaugeMetric{
			AbstractMetric: metrics.AbstractMetric{
				Name: metrics.TotalAlloc,
			},
			Value: metrics.Gauge(stats.TotalAlloc),
		},
	}
}
