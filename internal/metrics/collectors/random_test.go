package collectors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"eridiumdev/yandex-praktikum-go-devops/config"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

func TestRandomCollect(t *testing.T) {
	col, err := NewRandomCollector("random", config.RandomCollectorConfig{Min: 0, Max: 99})
	require.NoError(t, err)

	snapshot, err := col.Collect(context.Background())

	require.NoError(t, err)
	assert.Equal(t, 1, len(snapshot))
	assert.Equal(t, domain.RandomValue, snapshot[0].Name)
	assert.Equal(t, domain.TypeGauge, snapshot[0].Type)
	assert.GreaterOrEqual(t, snapshot[0].Gauge, domain.Gauge(0))
	assert.LessOrEqual(t, snapshot[0].Gauge, domain.Gauge(99))
}
