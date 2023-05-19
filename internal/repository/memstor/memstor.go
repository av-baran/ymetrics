package memstor

import (
	"fmt"
	"sync"

	"github.com/av-baran/ymetrics/internal/config"
	"github.com/av-baran/ymetrics/internal/metric"
	"github.com/av-baran/ymetrics/pkg/interrors"
)

type MemStorage struct {
	MetricsStor map[string]metric.Metric
}

var memStorageSync = sync.Mutex{}

func New() *MemStorage {
	return &MemStorage{
		MetricsStor: make(map[string]metric.Metric),
	}
}

func (s *MemStorage) SetMetric(m metric.Metric) error {
	memStorageSync.Lock()
	defer memStorageSync.Unlock()

	if m.MType != metric.GaugeType && m.MType != metric.CounterType {
		return interrors.ErrInvalidMetricType
	}

	var resultDelta, mDelta int64
	if m.Delta != nil {
		mDelta = *m.Delta
	}
	if s.MetricsStor[m.ID].Delta != nil {
		resultDelta = mDelta + *s.MetricsStor[m.ID].Delta
		m.Delta = &resultDelta
	}

	s.MetricsStor[m.ID] = m

	return nil
}

func (s *MemStorage) GetMetric(id string, mType string) (*metric.Metric, error) {
	memStorageSync.Lock()
	defer memStorageSync.Unlock()

	if mType != metric.GaugeType && mType != metric.CounterType {
		return nil, interrors.ErrInvalidMetricType
	}

	m, ok := s.MetricsStor[id]
	if !ok {
		return nil, interrors.ErrMetricNotFound
	}
	if m.MType != mType {
		return nil, interrors.ErrMetricNotFound
	}
	return &m, nil
}

func (s *MemStorage) GetAllMetrics() ([]metric.Metric, error) {
	memStorageSync.Lock()
	defer memStorageSync.Unlock()
	res := make([]metric.Metric, 0)

	for _, v := range s.MetricsStor {
		res = append(res, v)
	}

	return res, nil
}

func (s *MemStorage) SetMetricsBatch(metrics []metric.Metric) error {
	for _, m := range metrics {
		if err := s.SetMetric(m); err != nil {
			return fmt.Errorf("cannot update metrics with batch: %w", err)
		}
	}
	return nil
}

func (s *MemStorage) Init(cfg config.StorageConfig) error {
	return nil
}

func (s *MemStorage) Shutdown() error {
	return nil
}

func (s *MemStorage) Ping() error {
	return nil
}
