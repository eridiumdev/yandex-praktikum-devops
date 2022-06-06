package storage

import (
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestStoreMetric(t *testing.T) {
	mutex := &sync.RWMutex{}
	tests := []struct {
		name string
		have InMemStorage
		add  metrics.Metric
		want InMemStorage
	}{
		{
			name: "add counter to empty storage",
			have: InMemStorage{
				Metrics:      map[string]metrics.Metric{},
				metricsMutex: mutex,
			},
			add: metrics.CounterMetric{
				AbstractMetric: metrics.AbstractMetric{Name: metrics.PollCount},
				Value:          10,
			},
			want: InMemStorage{
				Metrics: map[string]metrics.Metric{
					metrics.PollCount: metrics.CounterMetric{
						AbstractMetric: metrics.AbstractMetric{Name: metrics.PollCount},
						Value:          10,
					},
				},
				metricsMutex: mutex,
			},
		},
		{
			name: "add gauge to empty storage",
			have: InMemStorage{
				Metrics:      map[string]metrics.Metric{},
				metricsMutex: mutex,
			},
			add: metrics.GaugeMetric{
				AbstractMetric: metrics.AbstractMetric{Name: metrics.Alloc},
				Value:          10.333,
			},
			want: InMemStorage{
				Metrics: map[string]metrics.Metric{
					metrics.Alloc: metrics.GaugeMetric{
						AbstractMetric: metrics.AbstractMetric{Name: metrics.Alloc},
						Value:          10.333,
					},
				},
				metricsMutex: mutex,
			},
		},
		{
			name: "update counter",
			have: InMemStorage{
				Metrics: map[string]metrics.Metric{
					metrics.PollCount: metrics.CounterMetric{
						AbstractMetric: metrics.AbstractMetric{Name: metrics.PollCount},
						Value:          10,
					},
				},
				metricsMutex: mutex,
			},
			add: metrics.CounterMetric{
				AbstractMetric: metrics.AbstractMetric{Name: metrics.PollCount},
				Value:          15,
			},
			want: InMemStorage{
				Metrics: map[string]metrics.Metric{
					metrics.PollCount: metrics.CounterMetric{
						AbstractMetric: metrics.AbstractMetric{Name: metrics.PollCount},
						Value:          25,
					},
				},
				metricsMutex: mutex,
			},
		},
		{
			name: "update gauge",
			have: InMemStorage{
				Metrics: map[string]metrics.Metric{
					metrics.Alloc: metrics.GaugeMetric{
						AbstractMetric: metrics.AbstractMetric{Name: metrics.Alloc},
						Value:          10.333,
					},
				},
				metricsMutex: mutex,
			},
			add: metrics.GaugeMetric{
				AbstractMetric: metrics.AbstractMetric{Name: metrics.Alloc},
				Value:          5.5,
			},
			want: InMemStorage{
				Metrics: map[string]metrics.Metric{
					metrics.Alloc: metrics.GaugeMetric{
						AbstractMetric: metrics.AbstractMetric{Name: metrics.Alloc},
						Value:          5.5,
					},
				},
				metricsMutex: mutex,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.have.StoreMetric(tt.add)
			require.NoError(t, err)
			assert.EqualValues(t, tt.want, tt.have)
		})
	}
}

func TestStoreMetricWithRaceCondition(t *testing.T) {
	store := &InMemStorage{
		Metrics:      map[string]metrics.Metric{},
		metricsMutex: &sync.RWMutex{},
	}
	metric := metrics.CounterMetric{
		AbstractMetric: metrics.AbstractMetric{Name: metrics.PollCount},
		Value:          1,
	}
	done := make(chan int)
	for i := 0; i < 1000; i++ {
		go func() {
			_ = store.StoreMetric(metric)
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
	assert.Equal(t, metrics.Counter(1000), store.Metrics[metrics.PollCount].GetValue())
}
