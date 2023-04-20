package server

import (
	"fmt"
	"net/http"

	"github.com/av-baran/ymetrics/internal/metric"
	"github.com/go-chi/chi/v5"
)

func (s *Server) GetMetricHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	name := chi.URLParam(r, "name")
	mType := metric.Type(chi.URLParam(r, "type"))

	switch mType {
	case metric.GaugeType:
		value, err := s.Storage.GetGauge(name)
		if err != nil {
			http.Error(w, fmt.Sprintf("cannot get gauge metric: %s", err), getErrorCode(err))
		}
		w.Write([]byte(value))
	case metric.CounterType:
		value, err := s.Storage.GetCounter(name)
		if err != nil {
			http.Error(w, fmt.Sprintf("cannot get gauge metric: %s", err), getErrorCode(err))
		}
		w.Write([]byte(value))
	default:
		http.Error(w, "unknown metric type", http.StatusBadRequest)
	}
}
