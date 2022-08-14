package service

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"eridiumdev/yandex-praktikum-go-devops/config"
	"eridiumdev/yandex-praktikum-go-devops/internal/common/logger"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

type metricsService struct {
	repo     MetricsRepository
	backuper MetricsBackuper
}

func NewMetricsService(
	ctx context.Context,
	repo MetricsRepository,
	backuper MetricsBackuper,
	backupCfg config.BackupConfig,
) (*metricsService, error) {
	s := &metricsService{
		repo:     repo,
		backuper: backuper,
	}
	if backuper != nil {
		if backupCfg.DoRestore {
			err := s.restoreFromLastBackup(ctx)
			if err != nil {
				return nil, errors.Wrap(err, "failed to restore from backup")
			}
		}
		if backupCfg.Interval > 0 {
			go s.startDoingBackups(ctx, backupCfg.Interval)
		}
	}
	return s, nil
}

func (s *metricsService) Update(ctx context.Context, metric domain.Metric) (domain.Metric, error) {
	existingMetric, found, err := s.repo.Get(ctx, metric.Name)
	if err != nil {
		return metric, err
	}
	if found && metric.IsCounter() {
		// For counters, old value is added on top of new value
		metric.Counter += existingMetric.Counter
	}
	return metric, s.repo.Store(ctx, metric)
}

func (s *metricsService) UpdateMany(ctx context.Context, metrics []domain.Metric) ([]domain.Metric, error) {
	names := make([]string, 0)
	for _, metric := range metrics {
		names = append(names, metric.Name)
	}

	existingMetrics, err := s.repo.List(ctx, &domain.MetricsFilter{Names: names})
	if err != nil {
		return metrics, err
	}

	for _, existingMetric := range existingMetrics {
		for i := range metrics {
			if existingMetric.IsCounter() {
				// For counters, old value is added on top of new value
				metrics[i].Counter += existingMetric.Counter
				break
			}
		}
	}
	return metrics, s.repo.Store(ctx, metrics...)
}

func (s *metricsService) Get(ctx context.Context, name string) (domain.Metric, bool, error) {
	return s.repo.Get(ctx, name)
}

func (s *metricsService) List(ctx context.Context) ([]domain.Metric, error) {
	return s.repo.List(ctx, nil)
}

func (s *metricsService) startDoingBackups(ctx context.Context, interval time.Duration) {
	backupCycles := 0
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			backupCycles++
			logger.New(ctx).Debugf("[metrics service] backup cycle %d begins", backupCycles)

			metrics, err := s.repo.List(ctx, nil)
			if err != nil {
				logger.New(ctx).Errorf("[metrics service] backup cycle %d failed, error: %s",
					backupCycles, err.Error())
				continue
			}

			err = s.backuper.Backup(metrics)
			if err != nil {
				logger.New(ctx).Errorf("[metrics service] backup cycle %d failed, error: %s",
					backupCycles, err.Error())
				continue
			}
			logger.New(ctx).Debugf("[metrics service] backup cycle %d successful, metrics count = %d",
				backupCycles, len(metrics))

		case <-ctx.Done():
			logger.New(ctx).Debugf("[metrics service] context cancelled, stopped doing backups")
			return
		}
	}
}

func (s *metricsService) restoreFromLastBackup(ctx context.Context) error {
	metrics, err := s.backuper.Restore()
	if err != nil {
		return err
	}
	logger.New(ctx).Infof("[metrics service] restored %d metrics from backup, applying to repo...", len(metrics))

	for _, metric := range metrics {
		err = s.repo.Store(ctx, metric)
		if err != nil {
			return err
		}
	}
	logger.New(ctx).Infof("[metrics service] %d metrics from backup restored to repository", len(metrics))
	return nil
}
