package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/av-baran/ymetrics/internal/metric"
	"github.com/go-chi/chi/v5"
)

func (s *Server) UpdateMetricHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	m := &metric.Metrics{
		ID:    chi.URLParam(r, "name"),
		MType: chi.URLParam(r, "type"),
	}
	value := chi.URLParam(r, "value")

	switch m.MType {
	case metric.GaugeType:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("cannot parse gauge metric: %s", err), http.StatusBadRequest)
		}
		m.Value = &v
		s.Storage.SetMetric(*m)
	case metric.CounterType:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("cannot parse counter metric: %s", err), http.StatusBadRequest)
		}
		m.Delta = &v
		s.Storage.AddCounter(m.ID, *m.Delta)
	default:
		http.Error(w, "unknown metric type", http.StatusNotImplemented)
		return
	}

	w.WriteHeader(http.StatusOK)
}
