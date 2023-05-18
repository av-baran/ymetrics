package psql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/av-baran/ymetrics/internal/config"
	"github.com/av-baran/ymetrics/internal/metric"
	"github.com/av-baran/ymetrics/pkg/interrors"
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

	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		return fmt.Errorf("cannot create new DB connection: %w", err)
	}
	s.db = db

	if err = s.Ping(); err != nil {
		return fmt.Errorf("cannot init DB: %w", err)
	}

	// ctx, cancel := context.WithTimeout(context.Background(), s.requestTimeout)
	// defer cancel()
	ctx := context.Background()
	err = interrors.RetryOnErr(
		func() error {
			return s.createTables(ctx)
		})
	if err != nil {
		return fmt.Errorf("cannot init DB: %w", err)
	}
	return nil
}

func (s *PsqlDB) SetMetric(m metric.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.requestTimeout)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("cannot begin tx: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
INSERT INTO metrics (id, type, value, delta) VALUES ($1, $2, $3, $4)
	ON CONFLICT (id) DO UPDATE SET
		value = EXCLUDED.value,
		type  = EXCLUDED.type, 
		delta = CASE WHEN metrics.delta IS NOT NULL OR EXCLUDED.delta IS NOT NULL
			THEN coalesce(metrics.delta, 0) + coalesce(EXCLUDED.delta, 0) 
			ELSE NULL END;`,
		m.ID, m.MType, m.Value, m.Delta,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("cannot exec query: %w", err)
	}
	return tx.Commit()
}

func (s *PsqlDB) GetMetric(id string, mType string) (*metric.Metric, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.requestTimeout)
	defer cancel()

	m := &metric.Metric{}

	row := s.db.QueryRowContext(ctx,
		`SELECT id, type, value, delta FROM metrics WHERE id=$1 AND type=$2`,
		id, mType)

	var val sql.NullFloat64
	var delta sql.NullInt64
	err := row.Scan(&m.ID, &m.MType, &val, &delta)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", interrors.ErrMetricNotFound, err)
	}
	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("cannot get all metrics: %w", err)
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
	res := make([]metric.Metric, 0)

	ctx := context.Background()
	rows, err := s.db.QueryContext(ctx, "SELECT id, type, value, delta FROM metrics")
	if err != nil {
		return nil, fmt.Errorf("cannot get all metrics: %w", err)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cannot get all metrics: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var val sql.NullFloat64
		var delta sql.NullInt64
		m := &metric.Metric{}

		err := rows.Scan(&m.ID, &m.MType, &val, &delta)
		if err != nil {
			return nil, fmt.Errorf("cannot get all metrics: %w", err)
		}
		if val.Valid {
			m.Value = &val.Float64
		}
		if delta.Valid {
			m.Delta = &delta.Int64
		}
		res = append(res, *m)
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

func (s *PsqlDB) UpdateBatch(metrics []metric.Metric) error {
	for _, m := range metrics {
		if err := s.SetMetric(m); err != nil {
			return fmt.Errorf("cannot update metrics with batch: %w", err)
		}
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
		return fmt.Errorf("%w: %w", interrors.ErrPingDB, err)
	}

	return nil
}

func (s *PsqlDB) createTables(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS metrics (
			"id" VARCHAR(256) PRIMARY KEY,
			"value"  DOUBLE PRECISION,
			"type" TEXT,
			"delta" BIGINT)`,
	)
	if err != nil {
		return fmt.Errorf("cannot create tables: %w", err)
	}

	return nil
}
