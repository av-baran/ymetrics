package repository

import (
	"github.com/av-baran/ymetrics/internal/config"
	"github.com/av-baran/ymetrics/internal/metric"
)

type Storage interface {
	Init(config.StorageConfig) error
	Ping() error
	Shutdown() error

	SetMetric(metric.Metric) error
	GetMetric(id string, mType string) (*metric.Metric, error)
	GetAllMetrics() ([]metric.Metric, error)
	SetMetricsBatch(m []metric.Metric) error
}
