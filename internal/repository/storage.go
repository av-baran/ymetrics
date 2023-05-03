package repository

import "github.com/av-baran/ymetrics/internal/metric"

type Storager interface {
	SetMetric(metric.Metrics) error
	GetMetric(*metric.Metrics) error
	GetAllMetrics() []metric.Metrics

	AddCounter(string, int64) int64
}
