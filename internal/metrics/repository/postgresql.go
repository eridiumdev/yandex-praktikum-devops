package repository

import (
	"context"
	"database/sql"

	// postgres init
	_ "github.com/jackc/pgx/v4/stdlib"

	"eridiumdev/yandex-praktikum-go-devops/config"
	"eridiumdev/yandex-praktikum-go-devops/internal/common/logger"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

type postgresRepo struct {
	db  *sql.DB
	cfg config.DatabaseConfig
}

func NewPostgresRepo(ctx context.Context, cfg config.DatabaseConfig) (*postgresRepo, error) {
	db, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		return nil, err
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()

	err = db.PingContext(timeoutCtx)
	if err != nil {
		return nil, err
	}

	return &postgresRepo{
		db:  db,
		cfg: cfg,
	}, nil
}

func (r *postgresRepo) Ping(ctx context.Context) bool {
	timeoutCtx, cancel := context.WithTimeout(ctx, r.cfg.ConnectTimeout)
	defer cancel()

	err := r.db.PingContext(timeoutCtx)
	if err != nil {
		logger.New(ctx).Errorf("[postgres repo] error on ping: %s", err.Error())
		return false
	}
	return true
}

func (r *postgresRepo) Store(metric domain.Metric) {
}

func (r *postgresRepo) Get(name string) (domain.Metric, bool) {
	return domain.Metric{}, true
}

func (r *postgresRepo) List() []domain.Metric {
	result := make([]domain.Metric, 0)
	return result
}
