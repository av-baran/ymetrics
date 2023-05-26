package repository

import (
	"context"

	"github.com/av-baran/ymetrics/internal/config"
	"github.com/av-baran/ymetrics/internal/logger"
	"github.com/av-baran/ymetrics/internal/metric"
	"github.com/av-baran/ymetrics/internal/repository/memstor"
	"github.com/av-baran/ymetrics/internal/repository/psql"
)

type Storage interface {
	Init(config.StorageConfig) error
	Ping(context.Context) error
	Shutdown() error

	SetMetric(context.Context, metric.Metric) error
	GetMetric(ctx context.Context, id string, mType string) (*metric.Metric, error)
	GetAllMetrics(context.Context) ([]metric.Metric, error)
	SetMetricsBatch(context.Context, []metric.Metric) error
}

func New(cfg config.StorageConfig) (Storage, error) {
	var repo Storage
	if cfg.DatabaseDSN != "" {
		repo = psql.New()
		err := repo.Init(cfg)
		if err != nil {
			logger.Fatalf("cannot init storage: %s", err)
		}
	} else {
		repo = memstor.New()
	}
	return repo, nil
}
