package repository

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

func TestStore(t *testing.T) {
	mutex := &sync.RWMutex{}
	tests := []struct {
		name string
		have inMemRepo
		add  domain.Metric
		want inMemRepo
	}{
		{
			name: "add counter to empty repo",
			have: inMemRepo{
				metrics: map[string]domain.Metric{},
				mutex:   mutex,
			},
			add: domain.NewCounter(domain.PollCount, 10),
			want: inMemRepo{
				metrics: map[string]domain.Metric{
					domain.PollCount: domain.NewCounter(domain.PollCount, 10),
				},
				mutex: mutex,
			},
		},
		{
			name: "add gauge to empty repo",
			have: inMemRepo{
				metrics: map[string]domain.Metric{},
				mutex:   mutex,
			},
			add: domain.NewGauge(domain.Alloc, 10.333),
			want: inMemRepo{
				metrics: map[string]domain.Metric{
					domain.Alloc: domain.NewGauge(domain.Alloc, 10.333),
				},
				mutex: mutex,
			},
		},
		{
			name: "update counter",
			have: inMemRepo{
				metrics: map[string]domain.Metric{
					domain.PollCount: domain.NewCounter(domain.PollCount, 10),
				},
				mutex: mutex,
			},
			add: domain.NewCounter(domain.PollCount, 15),
			want: inMemRepo{
				metrics: map[string]domain.Metric{
					domain.PollCount: domain.NewCounter(domain.PollCount, 25),
				},
				mutex: mutex,
			},
		},
		{
			name: "update gauge",
			have: inMemRepo{
				metrics: map[string]domain.Metric{
					domain.Alloc: domain.NewGauge(domain.Alloc, 10.333),
				},
				mutex: mutex,
			},
			add: domain.NewGauge(domain.Alloc, 5.5),
			want: inMemRepo{
				metrics: map[string]domain.Metric{
					domain.Alloc: domain.NewGauge(domain.Alloc, 5.5),
				},
				mutex: mutex,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.have.Store(tt.add)
			require.NoError(t, err)
			assert.EqualValues(t, tt.want, tt.have)
		})
	}
}

func TestStoreWithRaceCondition(t *testing.T) {
	repo := NewInMemRepo()
	metric := domain.NewCounter(domain.PollCount, 1)

	done := make(chan int)
	for i := 0; i < 1000; i++ {
		go func() {
			_ = repo.Store(metric)
			done <- 1
		}()
	}
	threadsDone := 0
	for range done {
		threadsDone++
		if threadsDone == 1000 {
			break
		}
	}
	result, err := repo.Get(metric.Name())
	require.NoError(t, err)
	assert.Equal(t, domain.Counter(1000), result.Value())
}
