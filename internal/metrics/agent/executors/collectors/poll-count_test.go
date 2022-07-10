package collectors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

func TestPollCountCollect(t *testing.T) {
	col := NewPollCountCollector("poll-count")
	snapshot, err := col.Collect(context.Background())

	require.NoError(t, err)
	assert.Equal(t, 1, len(snapshot))
	assert.Equal(t, domain.PollCount, snapshot[0].Name())
	assert.Equal(t, domain.TypeCounter, snapshot[0].Type())
	assert.Equal(t, domain.Counter(1), snapshot[0].Value())
}
