package service

import (
	"errors"
	"log"

	"github.com/av-baran/ymetrics/internal/entity/metric"
	"github.com/av-baran/ymetrics/pkg/interrors"
)

type Storager interface {
	GetMetricType(string) (metric.Type, bool)

	// FIXME хорошо ли здесь использовать пустой интерфейс или сделать отдельные методы под каждый тип
	StoreMetric(interface{}) error
	GetMetric(string) (interface{}, error)
}

type Service struct {
	Storage Storager
}

func New(s Storager) *Service {
	return &Service{Storage: s}
}

func (s *Service) UpdateGauge(m *metric.Gauge) error {
	t, ok := s.Storage.GetMetricType(m.Name)
	if ok && t != metric.GaugeType {
		return errors.New(interrors.ErrMetricExistsWithAnotherType)
	}

	if err := s.Storage.StoreMetric(m); err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateCounter(m *metric.Counter) error {
	t, ok := s.Storage.GetMetricType(m.Name)
	if ok && t != metric.CounterType {
		return errors.New(interrors.ErrMetricExistsWithAnotherType)
	}

	if ok {
		v, err := s.Storage.GetMetric(m.Name)
		if currentValue, ok := v.(int64); err == nil && ok {
			m.Value += currentValue
		} else {
			return errors.New(interrors.ErrStorageInternalError)
		}
	}
	if err := s.Storage.StoreMetric(m); err != nil {
		log.Printf("%v", err.Error())
		return err
	}
	return nil
}

func (s *Service) GetGauge(name string) (float64, error) {
	i, err := s.Storage.GetMetric(name)
	if err != nil {
		return 0.0, err
	}
	if v, ok := i.(float64); ok {
		return v, nil
	}
	return 0.0, errors.New(interrors.ErrStorageInternalError)
}

func (s *Service) GetCounter(name string) (int64, error) {
	i, err := s.Storage.GetMetric(name)
	if err != nil {
		return 0.0, err
	}
	if v, ok := i.(int64); ok {
		return v, nil
	}
	return 0, errors.New(interrors.ErrStorageInternalError)
}
