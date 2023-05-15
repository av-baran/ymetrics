package memstor

import (
	"sync"

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

	switch m.MType {
	case metric.GaugeType:
		if m.Value == nil {
			return interrors.ErrInvalidMetricValue
		}
		s.MetricsStor[m.ID] = m
	case metric.CounterType:
		if m.Delta == nil {
			return interrors.ErrInvalidMetricValue
		}
		if s.MetricsStor[m.ID].Delta != nil {
			currentValue := *s.MetricsStor[m.ID].Delta
			newValue := currentValue + *m.Delta
			m.Delta = &newValue
		}
		s.MetricsStor[m.ID] = m
	default:
		return interrors.ErrInvalidMetricType
	}
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

func (s *MemStorage) InitStorage(params string) error {
	return nil
}

func (s *MemStorage) Shutdown() error {
	return nil
}

func (s *MemStorage) Ping() error {
	return nil
}
