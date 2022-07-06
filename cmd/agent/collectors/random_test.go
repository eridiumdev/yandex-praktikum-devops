package collectors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"eridiumdev/yandex-praktikum-go-devops/internal/metrics"
)

func TestRandomCollect(t *testing.T) {
	col, err := NewRandomCollector("random", 0, 99)
	require.NoError(t, err)

	snapshot, err := col.Collect(context.Background())

	require.NoError(t, err)
	assert.Equal(t, 1, len(snapshot))
	assert.Equal(t, metrics.RandomValue, snapshot[0].GetName())
	assert.Equal(t, metrics.TypeGauge, snapshot[0].GetType())
	assert.GreaterOrEqual(t, snapshot[0].GetValue(), metrics.Gauge(0))
	assert.LessOrEqual(t, snapshot[0].GetValue(), metrics.Gauge(99))
}
