package memstor

import (
	"sync"

	"github.com/av-baran/ymetrics/internal/metric"
	"github.com/av-baran/ymetrics/pkg/interrors"
)

type MemStorage struct {
	GaugeStor   map[string]float64
	CounterStor map[string]int64
}

var memStorageSync = sync.Mutex{}

func New() *MemStorage {
	return &MemStorage{
		GaugeStor:   make(map[string]float64),
		CounterStor: make(map[string]int64),
	}
}

func (s *MemStorage) SetMetric(m metric.Metrics) error {
	memStorageSync.Lock()
	defer memStorageSync.Unlock()

	switch m.MType {
	case "gauge":
		if m.Value == nil {
			return interrors.ErrInvalidMetricValue
		}
		s.GaugeStor[m.ID] = *m.Value
	case "counter":
		if m.Delta == nil {
			return interrors.ErrInvalidMetricValue
		}
		s.CounterStor[m.ID] += *m.Delta
	default:
		return interrors.ErrInvalidMetricType
	}
	return nil
}

func (s *MemStorage) GetMetric(m *metric.Metrics) error {
	memStorageSync.Lock()
	defer memStorageSync.Unlock()

	switch m.MType {
	case "gauge":
		v, ok := s.GaugeStor[m.ID]
		if !ok {
			return interrors.ErrMetricNotFound
		}
		m.Value = &v
	case "counter":
		v, ok := s.CounterStor[m.ID]
		if !ok {
			return interrors.ErrMetricNotFound
		}
		m.Delta = &v
	default:
		return interrors.ErrInvalidMetricType
	}
	return nil
}

func (s *MemStorage) GetAllMetrics() []metric.Metrics {
	memStorageSync.Lock()
	defer memStorageSync.Unlock()
	res := make([]metric.Metrics, 0)

	for k, v := range s.GaugeStor {
		value := v
		m := &metric.Metrics{
			ID:    k,
			MType: "gauge",
			Delta: nil,
			Value: &value,
		}
		res = append(res, *m)
	}

	for k, v := range s.CounterStor {
		delta := v
		m := &metric.Metrics{
			ID:    k,
			MType: "counter",
			Delta: &delta,
			Value: nil,
		}
		res = append(res, *m)
	}
	return res
}

func (s *MemStorage) AddCounter(name string, value int64) int64 {
	memStorageSync.Lock()
	defer memStorageSync.Unlock()
	v := s.CounterStor[name]
	s.CounterStor[name] = v + value
	return s.CounterStor[name]
}
