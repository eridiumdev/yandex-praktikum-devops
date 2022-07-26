package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/repository"
)

func TestUpdate(t *testing.T) {
	tests := []struct {
		name    string
		updates []domain.Metric
		want    domain.Metric
	}{
		{
			name: "update counter one time",
			updates: []domain.Metric{
				domain.NewCounter(domain.PollCount, 10),
			},
			want: domain.NewCounter(domain.PollCount, 10),
		},
		{
			name: "update counter several times",
			updates: []domain.Metric{
				domain.NewCounter(domain.PollCount, 10),
				domain.NewCounter(domain.PollCount, 5),
				domain.NewCounter(domain.PollCount, 0),
			},
			want: domain.NewCounter(domain.PollCount, 15),
		},
		{
			name: "update gauge one time",
			updates: []domain.Metric{
				domain.NewGauge(domain.Alloc, 10.333),
			},
			want: domain.NewGauge(domain.Alloc, 10.333),
		},
		{
			name: "update gauge several times",
			updates: []domain.Metric{
				domain.NewGauge(domain.Alloc, 10.333),
				domain.NewGauge(domain.Alloc, 0.0),
				domain.NewGauge(domain.Alloc, 5.5),
			},
			want: domain.NewGauge(domain.Alloc, 5.5),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result domain.Metric
			s := NewMetricsService(repository.NewInMemRepo())

			for _, update := range tt.updates {
				result, _ = s.Update(update)
			}
			assert.Equal(t, tt.want, result)
		})
	}
}

func getDummyRepo() MetricsRepository {
	repo := repository.NewInMemRepo()
	repo.Store(domain.NewCounter(domain.PollCount, 10))
	repo.Store(domain.NewGauge(domain.Alloc, 10.333))
	return repo
}

func TestGet(t *testing.T) {
	type Want struct {
		metric domain.Metric
		found  bool
	}
	tests := []struct {
		name    string
		service *metricsService
		mName   string
		want    Want
	}{
		{
			name:    "get metric (found)",
			service: NewMetricsService(getDummyRepo()),
			mName:   domain.PollCount,
			want: Want{
				metric: domain.NewCounter(domain.PollCount, 10),
				found:  true,
			},
		},
		{
			name:    "get metric (not found)",
			service: NewMetricsService(getDummyRepo()),
			mName:   domain.RandomValue,
			want: Want{
				metric: domain.Metric{},
				found:  false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metric, found := tt.service.Get(tt.mName)
			assert.Equal(t, tt.want.metric, metric)
			assert.Equal(t, tt.want.found, found)
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
			list := tt.service.List()
			assert.ElementsMatch(t, tt.want, list)
		})
	}
}
