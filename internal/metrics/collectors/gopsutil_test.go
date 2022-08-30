package collectors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

func TestGopsutilCollect(t *testing.T) {
	col := NewGopsutilCollector("gopsutil")
	snapshot, err := col.Collect(context.Background())

	require.NoError(t, err)
	assert.Greater(t, len(snapshot), 0)

	expectedMetrics := []string{
		domain.TotalMemory,
		domain.FreeMemory,
		domain.CPUutilization1,
	}

	for _, m := range snapshot {
		if sliceContains(expectedMetrics, m.Name) {
			assert.Equal(t, domain.TypeGauge, m.Type)
			assert.NotNil(t, m.Gauge)
		}
	}
}

func sliceContains(slice []string, elem string) bool {
	for _, value := range slice {
		if value == elem {
			return true
		}
	}
	return false
}
