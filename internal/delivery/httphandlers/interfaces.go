package httphandlers

import "github.com/av-baran/ymetrics/internal/entity/metric"

//go:generate mockery --name=metricUpdater --exported=true
type metricUpdater interface {
	UpdateGauge(*metric.Gauge) error
	UpdateCounter(*metric.Counter) error
}

//go:generate mockery --name=metricGetter --exported=true
type metricGetter interface {
	GetGauge(string) (float64, error)
	GetCounter(string) (int64, error)
	GetAllMetrics()
}

//go:generate mockery --name Service
type Service interface {
	metricUpdater
	metricGetter
}
