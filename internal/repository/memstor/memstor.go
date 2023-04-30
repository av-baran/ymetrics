package memstor

import (
	"fmt"
	"sync"

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

func (s *MemStorage) SetGauge(name string, value float64) {
	memStorageSync.Lock()
	defer memStorageSync.Unlock()
	s.GaugeStor[name] = value
}

func (s *MemStorage) AddCounter(name string, value int64) {
	memStorageSync.Lock()
	defer memStorageSync.Unlock()
	v := s.CounterStor[name]
	s.CounterStor[name] = v + value
}

func (s *MemStorage) GetGauge(name string) (string, error) {
	memStorageSync.Lock()
	defer memStorageSync.Unlock()
	v, ok := s.GaugeStor[name]

	if !ok {
		return "", interrors.ErrMetricNotFound
	}
	return fmt.Sprintf("%v", v), nil
}

func (s *MemStorage) GetCounter(name string) (string, error) {
	memStorageSync.Lock()
	defer memStorageSync.Unlock()
	v, ok := s.CounterStor[name]
	if !ok {
		return "", interrors.ErrMetricNotFound
	}
	return fmt.Sprintf("%v", v), nil
}

func (s *MemStorage) GetAllGauge() map[string]float64 {
	memStorageSync.Lock()
	defer memStorageSync.Unlock()
	res := make(map[string]float64, len(s.GaugeStor))
	for k, v := range s.GaugeStor {
		res[k] = v
	}
	return res
}

func (s *MemStorage) GetAllCounter() map[string]int64 {
	memStorageSync.Lock()
	defer memStorageSync.Unlock()
	res := make(map[string]int64, len(s.GaugeStor))
	for k, v := range s.CounterStor {
		res[k] = v
	}
	return res
}
