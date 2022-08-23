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

const (
	StmtStoreMetrics = iota
	StmtGetMetrics
	StmtListMetrics
)

type postgresRepo struct {
	db    *sql.DB
	cfg   config.DatabaseConfig
	stmts map[int]*sql.Stmt
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
		m, migrateErr := migrate.New(fmt.Sprintf("file://%s", cfg.MigrationsDir), cfg.DSN)
		if migrateErr != nil {
			return nil, migrateErr
		}

		migrateErr = m.Up()
		switch {
		case migrateErr == nil:
			logger.New(ctx).Infof("[postgres repo] migrations successfully applied")
		case errors.Is(migrateErr, migrate.ErrNoChange):
			logger.New(ctx).Infof("[postgres repo] migrations: no change")
		default:
			return nil, migrateErr
		}
	}

	stmts, err := initStatements(ctx, db)
	if err != nil {
		return nil, err
	}

	return &postgresRepo{
		db:    db,
		cfg:   cfg,
		stmts: stmts,
	}, nil
}

func initStatements(ctx context.Context, db *sql.DB) (map[int]*sql.Stmt, error) {
	stmts := make(map[int]*sql.Stmt, 0)

	storeMetrics, err := db.PrepareContext(ctx,
		"INSERT INTO metrics (name, type, counter, gauge) VALUES ($1, $2, $3, $4)"+
			" ON CONFLICT (name) DO UPDATE SET counter = excluded.counter, gauge = excluded.gauge")
	if err != nil {
		return nil, err
	}
	stmts[StmtStoreMetrics] = storeMetrics

	getMetrics, err := db.PrepareContext(ctx, "SELECT name, type, counter, gauge FROM metrics WHERE name = $1")
	if err != nil {
		return nil, err
	}
	stmts[StmtGetMetrics] = getMetrics

	listMetrics, err := db.PrepareContext(ctx, "SELECT name, type, counter, gauge FROM metrics ORDER BY id desc")
	if err != nil {
		return nil, err
	}
	stmts[StmtListMetrics] = listMetrics

	return stmts, nil
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

func (r *postgresRepo) Store(ctx context.Context, metrics ...domain.Metric) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	stmt := tx.StmtContext(ctx, r.stmts[StmtStoreMetrics])

	for _, metric := range metrics {
		if _, err := stmt.ExecContext(ctx, metric.Name, metric.Type, metric.Counter, metric.Gauge); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *postgresRepo) Get(ctx context.Context, name string) (domain.Metric, bool, error) {
	var metric domain.Metric

	err := r.stmts[StmtGetMetrics].QueryRowContext(ctx, name).
		Scan(&metric.Name, &metric.Type, &metric.Counter, &metric.Gauge)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return metric, false, nil
	}
	return metric, err == nil, err
}

func (r *postgresRepo) List(ctx context.Context, filter *domain.MetricsFilter) ([]domain.Metric, error) {
	metrics := make([]domain.Metric, 0)

	var rows *sql.Rows
	var err error

	if filter != nil && len(filter.Names) > 0 {
		// Prepare args array (must have '[]any' type)
		args := make([]any, 0)
		// Build query with dynamic number of arguments
		query := "SELECT name, type, counter, gauge FROM metrics WHERE name IN ("
		for i, name := range filter.Names {
			if i > 0 {
				query += ","
			}
			query += fmt.Sprintf("$%d", i+1)
			args = append(args, name)
		}
		query += ")"
		rows, err = r.db.QueryContext(ctx, query, args...)
	} else {
		rows, err = r.stmts[StmtListMetrics].QueryContext(ctx)
	}

	if err != nil {
		return metrics, err
	}
	defer rows.Close()

	for rows.Next() {
		var metric domain.Metric
		err = rows.Scan(&metric.Name, &metric.Type, &metric.Counter, &metric.Gauge)
		if err != nil {
			return metrics, err
		}
		metrics = append(metrics, metric)
	}
	return metrics, rows.Err()
}
