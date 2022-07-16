package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/repository"
)

func TestUpdateCounter(t *testing.T) {
	type Update struct {
		name  string
		value domain.Counter
	}
	tests := []struct {
		name    string
		updates []Update
		want    domain.Metric
	}{
		{
			name: "update counter one time",
			updates: []Update{
				{
					domain.PollCount,
					domain.Counter(10),
				},
			},
			want: domain.NewCounter(domain.PollCount, 10),
		},
		{
			name: "update counter several times",
			updates: []Update{
				{
					domain.PollCount,
					domain.Counter(10),
				},
				{
					domain.PollCount,
					domain.Counter(20),
				},
				{
					domain.PollCount,
					domain.Counter(5),
				},
			},
			want: domain.NewCounter(domain.PollCount, 35),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result domain.Metric
			var err error
			s := NewMetricsService(repository.NewInMemRepo())

			for _, update := range tt.updates {
				result, err = s.UpdateCounter(update.name, update.value)
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestUpdateGauge(t *testing.T) {
	type Update struct {
		name  string
		value domain.Gauge
	}
	tests := []struct {
		name    string
		updates []Update
		want    domain.Metric
	}{
		{
			name: "update gauge one time",
			updates: []Update{
				{
					domain.Alloc,
					domain.Gauge(10.333),
				},
			},
			want: domain.NewGauge(domain.Alloc, 10.333),
		},
		{
			name: "update gauge several times",
			updates: []Update{
				{
					domain.Alloc,
					domain.Gauge(10.333),
				},
				{
					domain.Alloc,
					domain.Gauge(20.5),
				},
				{
					domain.Alloc,
					domain.Gauge(1.0),
				},
			},
			want: domain.NewGauge(domain.Alloc, 1.0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result domain.Metric
			var err error
			s := NewMetricsService(repository.NewInMemRepo())

			for _, update := range tt.updates {
				result, err = s.UpdateGauge(update.name, update.value)
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want, result)
		})
	}
}

func getDummyRepo() MetricsRepository {
	repo := repository.NewInMemRepo()
	_ = repo.Store(domain.NewCounter(domain.PollCount, 10))
	_ = repo.Store(domain.NewGauge(domain.Alloc, 10.333))
	return repo
}

func TestGet(t *testing.T) {
	tests := []struct {
		name    string
		service *metricsService
		mName   string
		want    domain.Metric
	}{
		{
			name:    "get metric (found)",
			service: NewMetricsService(getDummyRepo()),
			mName:   domain.PollCount,
			want:    domain.NewCounter(domain.PollCount, 10),
		},
		{
			name:    "get metric (not found)",
			service: NewMetricsService(getDummyRepo()),
			mName:   domain.RandomValue,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metric, err := tt.service.Get(tt.mName)
			require.NoError(t, err)
			assert.Equal(t, tt.want, metric)
		})
	}
}

func TestList(t *testing.T) {
	tests := []struct {
		name    string
		service *metricsService
		want    []domain.Metric
	}{
		{
			name:    "list metrics from service with empty repo",
			service: NewMetricsService(repository.NewInMemRepo()),
			want:    []domain.Metric{},
		},
		{
			name:    "list metrics from service with non-empty repo",
			service: NewMetricsService(getDummyRepo()),
			want: []domain.Metric{
				domain.NewCounter(domain.PollCount, 10),
				domain.NewGauge(domain.Alloc, 10.333),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := tt.service.List()
			require.NoError(t, err)
			assert.ElementsMatch(t, tt.want, list)
		})
	}
}
