package memstorv2

import (
	"fmt"
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

func (s *MemStorage) SetGauge(name string, value float64) {
	s.GaugeStor[name] = value
}

func (s *MemStorage) AddCounter(name string, value int64) {
	v, _ := s.CounterStor[name]
	s.CounterStor[name] = v + value
}

func (s *MemStorage) GetGauge(name string) string {
	return fmt.Sprintf("%v", s.GaugeStor[name])
}

func (s *MemStorage) GetCounter(name string) string {
	return fmt.Sprintf("%v", s.CounterStor[name])
}

func (s *MemStorage) GetAllGauge() map[string]float64 {
	return s.GaugeStor
}

func (s *MemStorage) GetAllCounter() map[string]int64 {
	return s.CounterStor
}
