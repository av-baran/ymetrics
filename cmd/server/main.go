package main

import (
	"errors"
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
	storage := &MemStorage{metrics: make(map[string]Metric)}

	mux := http.NewServeMux()
	mux.HandleFunc("/update/", storage.updateMetricHandler)
	mux.HandleFunc("/show/", storage.ShowHandler)
	return http.ListenAndServe("0.0.0.0:8080", mux)
}

type MemStorage struct {
	metrics map[string]Metric
}

func (s *MemStorage) updateMetricHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request with url: %v", r.URL.String())
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST request is allowed.", http.StatusMethodNotAllowed)
		return
	}

	metric, err := parsePath(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	s.updateMetric(metric)
	log.Printf("Metric updated. Current storage: %v", s)
}

func (s *MemStorage) updateMetric(m Metric) {
	log.Printf("Updating storage with metric: %v, value: %v", m.GetName(), m.GetValue())
	m.Update(s)
}

func (s *MemStorage) ShowHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Printing storage content: ")
	for k, v := range s.metrics {
		log.Printf("metrics[%v] -> %v = %v", k, v.GetName(), v.GetValue())
	}
}

func parsePath(path string) (Metric, error) {
	var metric Metric
	p := strings.Split(path, "/")
	if len(p) != 5 {
		return nil, errors.New("bad request: malformed URL")
	}

	rawType := p[2]
	rawName := p[3]
	rawValue := p[4]

	switch rawType {
	case "gauge":
		metric = &Gauge{}
	case "counter":
		metric = &Counter{}
	default:
		return nil, errors.New("bad request: invalid metric type")
	}
	if err := metric.SetName(rawName); err != nil {
		return nil, err
	}
	if err := metric.SetValue(rawValue); err != nil {
		return nil, err
	}
	return metric, nil
}

type Metric interface {
	SetName(string) error
	SetValue(string) error
	GetName() string
	GetValue() string
	Update(s *MemStorage)
}

type Gauge struct {
	Name  string
	Value float64
}

func (g *Gauge) SetName(s string) error {
	g.Name = s
	return nil
}

func (g *Gauge) SetValue(s string) error {
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return errors.New("bad request: wrong value")
	}
	g.Value = value
	return nil
}

func (g *Gauge) GetName() string {
	return g.Name
}

func (g *Gauge) GetValue() string {
	return strconv.FormatFloat(g.Value, 'e', 2, 64)
}

func (g *Gauge) Update(s *MemStorage) {
	s.metrics[g.Name] = g
}

type Counter struct {
	Name  string
	Value int64
}

func (c *Counter) SetName(s string) error {
	c.Name = s
	return nil
}

func (c *Counter) Update(s *MemStorage) {
	if _, ok := s.metrics[c.Name]; ok {
		storValue, err := strconv.ParseInt(s.metrics[c.Name].GetValue(), 10, 64)
		// FIXME
		if err != nil {
			return
		}
		c.Value += storValue
	}
	s.metrics[c.Name] = c
}

func (c *Counter) SetValue(s string) error {
	value, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return errors.New("bad request: wrong value")
	}
	c.Value = value
	return nil
}

func (c *Counter) GetName() string {
	return c.Name
}

func (c *Counter) GetValue() string {
	return strconv.FormatInt(c.Value, 10)
}
