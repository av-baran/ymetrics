package repository

import "github.com/av-baran/ymetrics/internal/metric"

type Storager interface {
	SetMetric(metric.Metrics) error

	SetGauge(string, float64)
	AddCounter(string, int64) int64

	GetGauge(string) (float64, error)
	GetCounter(string) (int64, error)

	GetAllMetrics() []metric.Metrics
}
