package collectors

import (
	"context"
	"errors"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"eridiumdev/yandex-praktikum-go-devops/internal/common/worker"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

var ErrEmptyCPUSnapshot = errors.New("could not scrape CPU metrics at this time (gopsutil returned empty slice)")

type gopsutilCollector struct {
	*worker.Worker
}

func NewGopsutilCollector(name string) *gopsutilCollector {
	col := &gopsutilCollector{
		Worker: worker.New(name, 1),
	}
	return col
}

func (col *gopsutilCollector) Collect(ctx context.Context) ([]domain.Metric, error) {
	memSnapshot, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return nil, err
	}

	cpuSnapshot, err := cpu.PercentWithContext(ctx, 0, false)
	if err != nil {
		return nil, err
	}
	if len(cpuSnapshot) < 1 {
		return nil, ErrEmptyCPUSnapshot
	}

	return []domain.Metric{
		domain.NewGauge(domain.TotalMemory, domain.Gauge(memSnapshot.Total)),
		domain.NewGauge(domain.FreeMemory, domain.Gauge(memSnapshot.Free)),
		domain.NewGauge(domain.CPUutilization1, domain.Gauge(cpuSnapshot[0])),
	}, nil
}
