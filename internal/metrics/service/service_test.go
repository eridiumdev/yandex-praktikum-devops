package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"eridiumdev/yandex-praktikum-go-devops/config"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/backup"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/repository"
)

func getDummyRepo() MetricsRepository {
	repo := repository.NewInMemRepo()
	ctx := context.Background()
	_ = repo.Store(ctx, domain.NewCounter(domain.PollCount, 10))
	_ = repo.Store(ctx, domain.NewGauge(domain.Alloc, 10.333))
	return repo
}

func getDummyBackuper() MetricsBackuper {
	return &backup.Mock{}
}

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
			ctx := context.Background()
			repo := repository.NewInMemRepo()
			backuper := getDummyBackuper()

			service, err := NewMetricsService(ctx, repo, backuper, config.BackupConfig{})
			require.NoError(t, err)

			var result domain.Metric

			for _, update := range tt.updates {
				result, err = service.Update(ctx, update)
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestGet(t *testing.T) {
	type Want struct {
		metric domain.Metric
		found  bool
	}
	tests := []struct {
		name  string
		mName string
		want  Want
	}{
		{
			name:  "get metric (found)",
			mName: domain.PollCount,
			want: Want{
				metric: domain.NewCounter(domain.PollCount, 10),
				found:  true,
			},
		},
		{
			name:  "get metric (not found)",
			mName: domain.RandomValue,
			want: Want{
				metric: domain.Metric{},
				found:  false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			repo := getDummyRepo()
			backuper := getDummyBackuper()

			service, err := NewMetricsService(ctx, repo, backuper, config.BackupConfig{})
			require.NoError(t, err)

			metric, found, err := service.Get(ctx, tt.mName)
			require.NoError(t, err)
			assert.Equal(t, tt.want.metric, metric)
			assert.Equal(t, tt.want.found, found)
		})
	}
}

func TestList(t *testing.T) {
	tests := []struct {
		name string
		repo MetricsRepository
		want []domain.Metric
	}{
		{
			name: "list metrics from service with empty repo",
			repo: repository.NewInMemRepo(),
			want: []domain.Metric{},
		},
		{
			name: "list metrics from service with non-empty repo",
			repo: getDummyRepo(),
			want: []domain.Metric{
				domain.NewCounter(domain.PollCount, 10),
				domain.NewGauge(domain.Alloc, 10.333),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			backuper := getDummyBackuper()

			service, err := NewMetricsService(ctx, tt.repo, backuper, config.BackupConfig{})
			require.NoError(t, err)

			list, err := service.List(ctx)
			require.NoError(t, err)
			assert.ElementsMatch(t, tt.want, list)
		})
	}
}
