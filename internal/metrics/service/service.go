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
	if backupCfg.DoRestore {
		err := s.restoreFromLastBackup(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to restore from backup")
		}
	}
	if backupCfg.Interval > 0 {
		go s.startDoingBackups(ctx, backupCfg.Interval)
	}
	return s, nil
}

func (s *metricsService) Update(metric domain.Metric) (updated domain.Metric, changed bool) {
	existingMetric, found := s.repo.Get(metric.Name)
	if found && metric.IsCounter() {
		// For counters, old value is added on top of new value
		metric.Counter += existingMetric.Counter
	}
	s.repo.Store(metric)
	return metric, metric != existingMetric
}

func (s *metricsService) Get(name string) (metric domain.Metric, found bool) {
	metric, found = s.repo.Get(name)
	return
}

func (s *metricsService) List() []domain.Metric {
	return s.repo.List()
}

func (s *metricsService) startDoingBackups(ctx context.Context, interval time.Duration) {
	backupCycles := 0
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			backupCycles++
			logger.New(ctx).Debugf("[metrics service] backup cycle %d begins", backupCycles)
			metrics := s.repo.List()
			if err := s.backuper.Backup(metrics); err == nil {
				logger.New(ctx).Debugf("[metrics service] backup cycle %d successful, metrics count = %d",
					backupCycles, len(metrics))
			}
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
	logger.New(ctx).Debugf("[metrics service] successfully restored %d metrics from backup", len(metrics))

	for _, metric := range metrics {
		s.repo.Store(metric)
	}
	return nil
}
