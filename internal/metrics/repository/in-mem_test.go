package repository

import (
	"context"
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
					domain.PollCount: domain.NewCounter(domain.PollCount, 15),
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
			ctx := context.Background()
			err := tt.have.Store(ctx, tt.add)
			require.NoError(t, err)
			assert.EqualValues(t, tt.want, tt.have)
		})
	}
}

func TestStoreWithRaceCondition(t *testing.T) {
	repo := NewInMemRepo()
	metric := domain.NewCounter(domain.PollCount, 1)

	done := make(chan int)
	ctx := context.Background()
	for i := 0; i < 1000; i++ {
		go func() {
			_ = repo.Store(ctx, metric)
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
	result, found, err := repo.Get(ctx, metric.Name)
	require.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, domain.Counter(1), result.Counter)
}

func TestGet(t *testing.T) {
	mutex := &sync.RWMutex{}
	type Want struct {
		metric domain.Metric
		found  bool
	}
	tests := []struct {
		name string
		repo *inMemRepo
		get  string
		want Want
	}{
		{
			name: "get metric from empty repo",
			repo: NewInMemRepo(),
			get:  domain.PollCount,
			want: Want{
				metric: domain.Metric{},
				found:  false,
			},
		},
		{
			name: "get metric from non-empty repo (found)",
			repo: &inMemRepo{
				metrics: map[string]domain.Metric{
					domain.Alloc:     domain.NewGauge(domain.Alloc, 10.333),
					domain.PollCount: domain.NewCounter(domain.PollCount, 10),
				},
				mutex: mutex,
			},
			get: domain.PollCount,
			want: Want{
				metric: domain.NewCounter(domain.PollCount, 10),
				found:  true,
			},
		},
		{
			name: "get metric from non-empty repo (not found)",
			repo: &inMemRepo{
				metrics: map[string]domain.Metric{
					domain.Alloc:     domain.NewGauge(domain.Alloc, 10.333),
					domain.PollCount: domain.NewCounter(domain.PollCount, 10),
				},
				mutex: mutex,
			},
			get: domain.HeapSys,
			want: Want{
				metric: domain.Metric{},
				found:  false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			metric, found, err := tt.repo.Get(ctx, tt.get)
			require.NoError(t, err)
			assert.Equal(t, tt.want.found, found)
			assert.Equal(t, tt.want.metric, metric)
		})
	}
}

func TestList(t *testing.T) {
	mutex := &sync.RWMutex{}
	tests := []struct {
		name string
		repo *inMemRepo
		want []domain.Metric
	}{
		{
			name: "get list from empty repo",
			repo: NewInMemRepo(),
			want: []domain.Metric{},
		},
		{
			name: "get list from non-empty repo",
			repo: &inMemRepo{
				metrics: map[string]domain.Metric{
					domain.Alloc:     domain.NewGauge(domain.Alloc, 10.333),
					domain.PollCount: domain.NewCounter(domain.PollCount, 10),
				},
				mutex: mutex,
			},
			want: []domain.Metric{
				domain.NewGauge(domain.Alloc, 10.333),
				domain.NewCounter(domain.PollCount, 10),
			},
		},
		{
			name: "get list from non-empty repo, different order",
			repo: &inMemRepo{
				metrics: map[string]domain.Metric{
					domain.PollCount: domain.NewCounter(domain.PollCount, 10),
					domain.Alloc:     domain.NewGauge(domain.Alloc, 10.333),
				},
				mutex: mutex,
			},
			want: []domain.Metric{
				domain.NewCounter(domain.PollCount, 10),
				domain.NewGauge(domain.Alloc, 10.333),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			list, err := tt.repo.List(ctx)
			require.NoError(t, err)
			assert.ElementsMatch(t, tt.want, list)
		})
	}
}
