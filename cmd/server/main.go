package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	storage := NewMemStorage()

	mux := http.NewServeMux()
	mux.HandleFunc("/update/", storage.updateMetricsHandler)
	mux.HandleFunc("/show/", storage.showAllHandler)
	return http.ListenAndServe("0.0.0.0:8080", mux)
}

type MemStorage struct {
	// FIXME Нужен ли дополнительный словарь чтобы хранить название и тип существующих метрик? Или просто искать имя в обеих мапах?
	metrics     map[string]string
	gaugeStor   map[string]float64
	counterStor map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics:     make(map[string]string),
		gaugeStor:   make(map[string]float64),
		counterStor: make(map[string]int64),
	}
}

// FIXME Привязывать ли хендлер к структуре или сделать вызов функции которая обернет логику в хендлер?
func (s *MemStorage) updateMetricsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request with url: %v", r.URL.String())
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST request is allowed.", http.StatusMethodNotAllowed)
		return
	}
	metric, err := parseURL(r.URL.Path)
	if err != nil {
		log.Printf("Error while parsing url")
		http.Error(w, err.Msg, err.Code)
		return
	}
	if err := s.updateMetric(metric); err != nil {
		log.Printf("Error while updating metric")
		http.Error(w, err.Msg, err.Code)
		return
	}
}

func (s *MemStorage) updateMetric(m *RawMetric) *httpError {
	log.Printf("Updating storage with metric: %v, type: %v, value: %v", m.Name, m.mType, m.Value)
	switch m.mType {
	case "gauge":
		return s.updateGauge(m)
	case "counter":
		return s.updateCounter(m)
	default:
		return &httpError{"invalid metric type", http.StatusNotImplemented}
	}
}

func (s *MemStorage) updateGauge(m *RawMetric) *httpError {
	parsedValue, err := strconv.ParseFloat(m.Value, 64)
	if err != nil {
		return &httpError{"invalid value", http.StatusBadRequest}
	}
	log.Printf("Parsed value of %v is %v", m.Value, parsedValue)

	if err := s.addMetric(m); err != nil {
		return err
	}

	log.Printf("Storing new value")
	s.gaugeStor[m.Name] = parsedValue
	log.Printf("%v is stored in gauge storage. Current values is %v", m.Name, s.gaugeStor[m.Name])
	return nil
}

func (s *MemStorage) updateCounter(m *RawMetric) *httpError {
	parsedValue, err := strconv.ParseInt(m.Value, 10, 64)
	if err != nil {
		return &httpError{"invalid value", http.StatusInternalServerError}
	}
	log.Printf("Parsed value of %v is %v", m.Value, parsedValue)

	if err := s.addMetric(m); err != nil {
		return err
	}

	s.counterStor[m.Name] += parsedValue
	log.Printf("%v is stored in counter storage. Current value is %v", m.Name, s.counterStor[m.Name])
	return nil
}

func (s *MemStorage) addMetric(m *RawMetric) *httpError {
	log.Printf("Checking if metrics exists")
	existingMetric, ok := s.metrics[m.Name]
	if ok && existingMetric != m.mType {
		return &httpError{"metric with same name and different type already exists", http.StatusBadRequest}
	} else if !ok {
		s.metrics[m.Name] = m.mType
	}
	return nil
}

func (s *MemStorage) showAllHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Registred metrics: ")
	for k, v := range s.metrics {
		log.Printf("%v type is %v", k, v)
		switch v {
		case "gauge":
			log.Printf("Value is %v", s.gaugeStor[k])
		case "counter":
			log.Printf("Value is %v", s.counterStor[k])
		}
	}
}

// FIXME Как лучше вернуть три значения? Выносить ли в отдельную функцию? Привязывать ли функцию к структуре, она вроде связанна логически, но данные структуры не нужны?
func parseURL(path string) (*RawMetric, *httpError) {
	p := strings.Split(path, "/")
	if len(p) != 5 {
		return nil, &httpError{"bad request: malformed URL", http.StatusNotFound}
	}
	return &RawMetric{
		mType: p[2],
		Name:  p[3],
		Value: p[4],
	}, nil
}

type RawMetric struct {
	mType string
	Name  string
	Value string
}

type httpError struct {
	Msg  string
	Code int
}

func (e *httpError) Error() string {
	return fmt.Sprintf("Code: %v, Message: %v", e.Code, e.Msg)
}
