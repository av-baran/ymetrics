package repository

import "github.com/av-baran/ymetrics/internal/metric"

type Storager interface {
	SetMetric(metric.Metric) error
	GetMetric(id string, mType string) (*metric.Metric, error)
	GetAllMetrics() []metric.Metric
}
