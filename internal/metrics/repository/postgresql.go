package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	// postgres init
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pkg/errors"

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
	logger.New(ctx).Infof("[postgres repo] connected to database")

	if cfg.MigrationsDir != "" {
		// Run migrations
		m, err := migrate.New(fmt.Sprintf("file://%s", cfg.MigrationsDir), cfg.DSN)
		if err != nil {
			return nil, err
		}

		migrateErr := m.Up()
		switch {
		case migrateErr == nil:
			logger.New(ctx).Infof("[postgres repo] migrations successfully applied")
		case errors.Is(migrateErr, migrate.ErrNoChange):
			logger.New(ctx).Infof("[postgres repo] migrations: no change")
		default:
			return nil, migrateErr
		}
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

func (r *postgresRepo) Store(ctx context.Context, metric domain.Metric) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO metrics (name, type, counter, gauge) VALUES ($1, $2, $3, $4)"+
		"ON CONFLICT (name) DO UPDATE SET counter = excluded.counter, gauge = excluded.gauge",
		metric.Name, metric.Type, metric.Counter, metric.Gauge)

	return err
}

func (r *postgresRepo) Update(ctx context.Context, metric domain.Metric) error {
	_, err := r.db.ExecContext(ctx, "UPDATE metrics SET counter = $1, gauge = $2 WHERE name = $3",
		metric.Counter, metric.Gauge, metric.Name)

	return err
}

func (r *postgresRepo) Get(ctx context.Context, name string) (domain.Metric, bool, error) {
	var metric domain.Metric

	err := r.db.QueryRowContext(ctx, "SELECT name, type, counter, gauge FROM metrics WHERE name = $1", name).
		Scan(&metric.Name, &metric.Type, &metric.Counter, &metric.Gauge)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return metric, false, nil
	}
	return metric, err == nil, err
}

func (r *postgresRepo) List(ctx context.Context) ([]domain.Metric, error) {
	metrics := make([]domain.Metric, 0)

	rows, err := r.db.QueryContext(ctx, "SELECT * FROM metrics ORDER BY id desc")
	if err != nil {
		return metrics, err
	}
	defer rows.Close()

	for rows.Next() {
		var metric domain.Metric
		err = rows.Scan(&metric)
		if err != nil {
			return metrics, err
		}
	}
	return metrics, rows.Err()
}
