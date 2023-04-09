package mem

import (
	"errors"
	"log"
	"strconv"

	"github.com/av-baran/ymetrics/internal/entities/metric"
	"github.com/av-baran/ymetrics/internal/interrors"
)

type MemStorage struct {
	// FIXME Нужен ли дополнительный словарь чтобы хранить название и тип существующих метрик? Или просто искать имя в обеих мапах?
	metrics     map[string]metric.Type
	gaugeStor   map[string]float64
	counterStor map[string]int64
}

func New() *MemStorage {
	return &MemStorage{
		metrics:     make(map[string]metric.Type),
		gaugeStor:   make(map[string]float64),
		counterStor: make(map[string]int64),
	}
}

func (s *MemStorage) UpdateMetric(m *metric.Rawdata) error {
	log.Printf("Updating storage with metric: %v, type: %v, value: %v", m.Name, m.Type, m.Value)
	switch m.Type {
	case metric.Gauge:
		return s.updateGauge(m)
	case metric.Counter:
		return s.updateCounter(m)
	default:
		return errors.New(interrors.ErrInvalidMetricType)
	}
}

func (s *MemStorage) updateGauge(m *metric.Rawdata) error {
	parsedValue, err := strconv.ParseFloat(m.Value, 64)
	if err != nil {
		return errors.New(interrors.ErrInvalidMetricValue)
	}

	if err := s.addMetric(m); err != nil {
		return err
	}

	s.gaugeStor[m.Name] = parsedValue
	log.Printf("%v is stored in gauge storage. Current values is %v", m.Name, s.gaugeStor[m.Name])
	return nil
}

func (s *MemStorage) updateCounter(m *metric.Rawdata) error {
	parsedValue, err := strconv.ParseInt(m.Value, 10, 64)
	if err != nil {
		return errors.New(interrors.ErrInvalidMetricValue)
	}
	log.Printf("Parsed value of %v is %v", m.Value, parsedValue)

	if err := s.addMetric(m); err != nil {
		return err
	}

	s.counterStor[m.Name] += parsedValue
	log.Printf("%v is stored in counter storage. Current value is %v", m.Name, s.counterStor[m.Name])
	return nil
}

func (s *MemStorage) addMetric(m *metric.Rawdata) error {
	existingMetric, ok := s.metrics[m.Name]
	if ok && existingMetric != m.Type {
		return errors.New(interrors.ErrMetricAlreadyExists)
	} else if !ok {
		s.metrics[m.Name] = m.Type
	}
	return nil
}
