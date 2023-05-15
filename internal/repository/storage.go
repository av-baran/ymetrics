package repository

import "github.com/av-baran/ymetrics/internal/metric"

type Storage interface {
	InitStorage(params string) error
	Ping() error
	Shutdown() error

	SetMetric(metric.Metric) error
	GetMetric(id string, mType string) (*metric.Metric, error)
	GetAllMetrics() ([]metric.Metric, error)
	UpdateBatch(m []metric.Metric) error
}
