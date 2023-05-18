package psql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/av-baran/ymetrics/internal/config"
	"github.com/av-baran/ymetrics/internal/metric"
	"github.com/av-baran/ymetrics/pkg/interrors"
)

const (
	createTableStatement = `
CREATE TABLE IF NOT EXISTS metrics (
	"id" VARCHAR(256) PRIMARY KEY,
	"value"  DOUBLE PRECISION,
	"type" TEXT,
	"delta" BIGINT);`

	setStatement = `
INSERT INTO metrics (id, type, value, delta) VALUES ($1, $2, $3, $4)
	ON CONFLICT (id) DO UPDATE SET
		value = EXCLUDED.value,
		type  = EXCLUDED.type,
		delta = CASE WHEN metrics.delta IS NOT NULL OR EXCLUDED.delta IS NOT NULL
			THEN coalesce(metrics.delta, 0) + coalesce(EXCLUDED.delta, 0)
			ELSE NULL END;`

	getStatement    = `SELECT id, type, value, delta FROM metrics WHERE id=$1 AND type=$2 LIMIT 1;`
	getAllStatement = `SELECT id, type, value, delta FROM metrics LIMIT 1000;`
)

type PsqlDB struct {
	db             *sql.DB
	requestTimeout time.Duration
}

func New() *PsqlDB {
	return &PsqlDB{}
}

func (s *PsqlDB) Init(cfg config.StorageConfig) error {
	s.requestTimeout = cfg.RequestTimeout

	ctx, cancel := context.WithTimeout(context.Background(), s.requestTimeout)
	defer cancel()

	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		return fmt.Errorf("cannot create new DB connection: %w", err)
	}
	s.db = db

	if err = s.Ping(); err != nil {
		return fmt.Errorf("DB is not available: %w", err)
	}

	stmt, err := s.db.PrepareContext(ctx, createTableStatement)
	if err != nil {
		return fmt.Errorf("cannot prepare create table statements: %w", err)
	}
	defer stmt.Close()

	err = interrors.RetryOnErr(func() error {
		_, err = stmt.ExecContext(ctx)
		return err
	})
	if err != nil {
		return fmt.Errorf("cannot init metric table: %w", err)
	}

	return nil
}

func (s *PsqlDB) SetMetric(m metric.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.requestTimeout)
	defer cancel()

	stmt, err := s.db.PrepareContext(ctx, setStatement)
	if err != nil {
		return fmt.Errorf("cannot prepare set metric statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, m.ID, m.MType, m.Value, m.Delta)
	if err != nil {
		return fmt.Errorf("cannot exec statement: %w", err)
	}

	return nil
}

func (s *PsqlDB) GetMetric(id string, mType string) (*metric.Metric, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.requestTimeout)
	defer cancel()

	m := &metric.Metric{}

	stmt, err := s.db.PrepareContext(ctx, getStatement)
	if err != nil {
		return nil, fmt.Errorf("cannot get metric: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id, mType)

	var val sql.NullFloat64
	var delta sql.NullInt64
	err = row.Scan(&m.ID, &m.MType, &val, &delta)
	if err != nil {
		resErr := errors.Join(interrors.ErrMetricNotFound, err)
		return nil, resErr
	}

	if val.Valid {
		m.Value = &val.Float64
	}
	if delta.Valid {
		m.Delta = &delta.Int64
	}

	return m, nil
}

func (s *PsqlDB) GetAllMetrics() ([]metric.Metric, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.requestTimeout)
	defer cancel()

	res := make([]metric.Metric, 0)

	stmt, err := s.db.PrepareContext(ctx, getAllStatement)
	if err != nil {
		return nil, fmt.Errorf("cannot prepare get all metrics statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot query all metrics: %w", err)
	}
	// Без этой проверки не проходит statictest, но нужна ли она? Т.к. мы проверили результат выполнения QueryContext
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("query result contains error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var val sql.NullFloat64
		var delta sql.NullInt64
		m := metric.Metric{}

		err := rows.Scan(&m.ID, &m.MType, &val, &delta)
		if err != nil {
			return nil, fmt.Errorf("cannot scan query result: %w", err)
		}

		if val.Valid {
			m.Value = &val.Float64
		}
		if delta.Valid {
			m.Delta = &delta.Int64
		}

		res = append(res, m)
	}

	return res, nil
}

func (s *PsqlDB) Shutdown() error {
	if s.db != nil {
		err := s.db.Close()
		if err != nil {
			return fmt.Errorf("cannot close DB connection: %w", err)
		}
	}

	return nil
}

func (s *PsqlDB) SetMetricsBatch(metrics []metric.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.requestTimeout)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("cannot begin transaction: %w", err)
	}

	stmt, err := tx.PrepareContext(ctx, setStatement)
	if err != nil {
		return fmt.Errorf("cannot prepare set metric statement: %w", err)
	}
	defer stmt.Close()

	for _, m := range metrics {
		_, err = tx.StmtContext(ctx, stmt).Exec(m.ID, m.MType, m.Value, m.Delta)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("cannot exec statement: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("cannot commit transaction: %w", err)
	}

	return nil
}

func (s *PsqlDB) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := interrors.RetryOnErr(
		func() error {
			return s.db.PingContext(ctx)
		})
	if err != nil {
		resError := errors.Join(interrors.ErrPingDB, err)
		return resError
	}

	return nil
}
