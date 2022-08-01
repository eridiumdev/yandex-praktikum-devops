package backup

import (
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

type Mock struct{}

func (b *Mock) Backup(metrics []domain.Metric) error {
	return nil
}

func (b *Mock) Restore() ([]domain.Metric, error) {
	return []domain.Metric{}, nil
}
