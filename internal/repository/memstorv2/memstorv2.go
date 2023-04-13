package memstorv2

import (
	"errors"
	"log"

	"github.com/av-baran/ymetrics/internal/entity/metric"
	"github.com/av-baran/ymetrics/pkg/interrors"
)

type MemStorage struct {
	gaugeStor   map[string]float64
	counterStor map[string]int64
}

func New() *MemStorage {
	return &MemStorage{
		gaugeStor:   make(map[string]float64),
		counterStor: make(map[string]int64),
	}
}

func (s *MemStorage) StoreMetric(m interface{}) error {
	switch m := m.(type) {
	case *metric.Gauge:
		s.gaugeStor[m.Name] = m.Value
		log.Printf("%v stored in gauge storage. Current values is %v", m.Name, s.gaugeStor[m.Name])
	case *metric.Counter:
		s.counterStor[m.Name] = m.Value
		log.Printf("%v stored in counter storage. Current value is %v", m.Name, s.counterStor[m.Name])
	default:
		return errors.New(interrors.ErrInvalidMetricType)
	}
	return nil
}

// FIXME
func (s *MemStorage) GetMetricType(name string) (metric.Type, bool) {
	if _, ok := s.gaugeStor[name]; ok {
		return metric.GaugeType, true
	}
	if _, ok := s.counterStor[name]; ok {
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
		return s.gaugeStor[name], nil
	case metric.CounterType:
		return s.counterStor[name], nil
	default:
		return nil, errors.New(interrors.ErrInvalidMetricType)
	}
}
