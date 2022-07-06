package collectors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"eridiumdev/yandex-praktikum-go-devops/internal/metrics"
)

func TestPollCountCollect(t *testing.T) {
	col := NewPollCountCollector("poll-count")
	snapshot, err := col.Collect(context.Background())

	require.NoError(t, err)
	assert.Equal(t, 1, len(snapshot))
	assert.Equal(t, metrics.PollCount, snapshot[0].GetName())
	assert.Equal(t, metrics.TypeCounter, snapshot[0].GetType())
	assert.Equal(t, metrics.Counter(1), snapshot[0].GetValue())
}
