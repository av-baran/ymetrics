package psql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/av-baran/ymetrics/internal/config"
	"github.com/av-baran/ymetrics/internal/metric"
	"github.com/av-baran/ymetrics/pkg/interrors"
)

type PsqlDB struct {
	db                *sql.DB
	queryTimeout      time.Duration
	batchQueryTimeout time.Duration
}

func New() *PsqlDB {
	return &PsqlDB{}
}

func (s *PsqlDB) Init(cfg config.StorageConfig) error {
	s.queryTimeout = cfg.QueryTimeout
	s.batchQueryTimeout = s.queryTimeout * 10

	initCtx, cancel := context.WithTimeout(context.Background(), s.queryTimeout)
	defer cancel()

	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		return fmt.Errorf("cannot create new DB connection: %w", err)
	}
	s.db = db

	if err = s.Ping(initCtx); err != nil {
		return fmt.Errorf("DB is not available: %w", err)
	}

	if err = s.applyMigrations(initCtx); err != nil {
		return fmt.Errorf("cannot apply migrations: %w", err)
	}

	return nil
}

func (s *PsqlDB) SetMetric(ctx context.Context, m metric.Metric) error {
	queryCtx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	stmtString := fmt.Sprintf(setStatement, "($1, $2, $3, $4)")

	stmt, err := s.db.PrepareContext(queryCtx, stmtString)
	if err != nil {
		return fmt.Errorf("cannot prepare set metric statement: %w", err)
	}
	defer stmt.Close()

	err = interrors.RetryOnErr(func() error {
		_, err = stmt.ExecContext(queryCtx, m.ID, m.MType, m.Value, m.Delta)
		return err
	})
	if err != nil {
		return fmt.Errorf("cannot exec statement: %w", err)
	}

	return nil
}

func (s *PsqlDB) GetMetric(ctx context.Context, id string, mType string) (*metric.Metric, error) {
	queryCtx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	m := &metric.Metric{}

	stmt, err := s.db.PrepareContext(queryCtx, getStatement)
	if err != nil {
		return nil, fmt.Errorf("cannot get metric: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(queryCtx, id, mType)

	var val sql.NullFloat64
	var delta sql.NullInt64
	err = row.Scan(&m.ID, &m.MType, &val, &delta)
	if err != nil {
		resErr := errors.Join(interrors.ErrMetricNotFound, err)
		return nil, fmt.Errorf("cannot get metric: %w", resErr)
	}

	if val.Valid {
		m.Value = &val.Float64
	}
	if delta.Valid {
		m.Delta = &delta.Int64
	}

	return m, nil
}

func (s *PsqlDB) GetAllMetrics(ctx context.Context) ([]metric.Metric, error) {
	queryCtx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	res := make([]metric.Metric, 0)

	stmt, err := s.db.PrepareContext(queryCtx, getAllStatement)
	if err != nil {
		return nil, fmt.Errorf("cannot prepare get all metrics statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(queryCtx)
	if err != nil {
		return nil, fmt.Errorf("cannot query all metrics: %w", err)
	}
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

func (s *PsqlDB) SetMetricsBatch(ctx context.Context, metrics []metric.Metric) error {
	metrics = deduplicateMetrics(metrics)

	queryValues := make([]string, 0, len(metrics))
	queryArgs := make([]interface{}, 0, len(metrics)*4)
	for i, m := range metrics {
		queryValues = append(queryValues, fmt.Sprintf(
			"($%d, $%d, $%d, $%d)",
			i*4+1, i*4+2, i*4+3, i*4+4,
		))
		queryArgs = append(queryArgs, m.ID)
		queryArgs = append(queryArgs, m.MType)
		queryArgs = append(queryArgs, m.Value)
		queryArgs = append(queryArgs, m.Delta)
	}

	stmtString := fmt.Sprintf(setStatement, strings.Join(queryValues, ","))

	queryCtx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	tx, err := s.db.BeginTx(queryCtx, nil)
	if err != nil {
		return fmt.Errorf("cannot begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(queryCtx, stmtString)
	if err != nil {
		return fmt.Errorf("cannot prepare set metric statement: %w", err)
	}
	defer stmt.Close()

	err = interrors.RetryOnErr(func() error {
		_, err = tx.StmtContext(queryCtx, stmt).Exec(queryArgs...)
		return err
	})
	if err != nil {
		return fmt.Errorf("cannot exec statement: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("cannot commit transaction: %w", err)
	}

	return nil
}

func (s *PsqlDB) Ping(ctx context.Context) error {
	pingCtx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	err := interrors.RetryOnErr(func() error {
		return s.db.PingContext(pingCtx)
	})
	if err != nil {
		resError := errors.Join(interrors.ErrPingDB, err)
		return fmt.Errorf("cannot ping db: %w", resError)
	}

	return nil
}

func (s *PsqlDB) applyMigrations(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.queryTimeout)
	defer cancel()

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

func deduplicateMetrics(metrics []metric.Metric) []metric.Metric {
	var dedupedMetrics []metric.Metric

	metricsMap := make(map[string]metric.Metric, 0)
	for _, m := range metrics {
		var resultDelta, mDelta int64
		if m.Delta != nil {
			mDelta = *m.Delta
		}
		if metricsMap[m.ID].Delta != nil {
			resultDelta = mDelta + *metricsMap[m.ID].Delta
			m.Delta = &resultDelta
		}

		metricsMap[m.ID] = m
	}

	for _, m := range metricsMap {
		dedupedMetrics = append(dedupedMetrics, m)
	}

	return dedupedMetrics
}
