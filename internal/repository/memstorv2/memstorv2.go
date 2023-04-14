package memstorv2

import (
	"errors"
	"log"

	"github.com/av-baran/ymetrics/internal/entity/metric"
	"github.com/av-baran/ymetrics/pkg/interrors"
)

type MemStorage struct {
	GaugeStor   map[string]float64
	CounterStor map[string]int64
}

func New() *MemStorage {
	return &MemStorage{
		GaugeStor:   make(map[string]float64),
		CounterStor: make(map[string]int64),
	}
}

func (s *MemStorage) StoreMetric(m interface{}) error {
	switch m := m.(type) {
	case *metric.Gauge:
		s.GaugeStor[m.Name] = m.Value
		log.Printf("%v stored in gauge storage. Current values is %v", m.Name, s.GaugeStor[m.Name])
	case *metric.Counter:
		s.CounterStor[m.Name] = m.Value
		log.Printf("%v stored in counter storage. Current value is %v", m.Name, s.CounterStor[m.Name])
	default:
		return errors.New(interrors.ErrInvalidMetricType)
	}
	return nil
}

// FIXME
func (s *MemStorage) GetMetricType(name string) (metric.Type, bool) {
	if _, ok := s.GaugeStor[name]; ok {
		return metric.GaugeType, true
	}
	if _, ok := s.CounterStor[name]; ok {
		return metric.CounterType, true
	}
	return metric.UnknownType, false
}

func (s *MemStorage) GetMetric(name string) (interface{}, error) {
	mType, ok := s.GetMetricType(name)
	if !ok {
		return nil, errors.New(interrors.ErrMetricNotFound)
	}
	switch mType {
	case metric.GaugeType:
		return s.GaugeStor[name], nil
	case metric.CounterType:
		return s.CounterStor[name], nil
	default:
		return nil, errors.New(interrors.ErrInvalidMetricType)
	}
}

func (s *MemStorage) GetAllGauge() map[string]float64 {
	return s.GaugeStor
}

func (s *MemStorage) GetAllCounter() map[string]int64 {
	return s.CounterStor
}
