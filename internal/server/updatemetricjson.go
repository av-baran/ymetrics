package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/av-baran/ymetrics/internal/metric"
)

func (s *Server) UpdateMetricJSONHandler(w http.ResponseWriter, r *http.Request) {
	var m metric.Metrics

	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		err = fmt.Errorf("cannot decode body: %w", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("update body: %+v", m)

	switch m.MType {
	case string(metric.CounterType):
		s.Storage.AddCounter(m.ID, *m.Delta)

		// FIXME
		stringVal, err := s.Storage.GetCounter(m.ID)
		if err != nil {
			err = fmt.Errorf("cannot get counter: %w", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if v, err := strconv.Atoi(stringVal); err != nil {
			err = fmt.Errorf("cannot convert counter to int: %w", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			*m.Delta = int64(v)
		}

	case string(metric.GaugeType):
		s.Storage.SetGauge(m.ID, *m.Value)

		// FIXME
		stringVal, err := s.Storage.GetGauge(m.ID)
		if err != nil {
			err = fmt.Errorf("cannot get gauge: %w", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if v, err := strconv.ParseFloat(stringVal, 64); err != nil {
			err = fmt.Errorf("cannot convert gauge to float64: %w", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			*m.Value = v
		}
	default:
		http.Error(w, "unknown metric type", http.StatusNotImplemented)
		return
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&m); err != nil {
		err = fmt.Errorf("cannot encode metric: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("get resp body: %+v", m)
}