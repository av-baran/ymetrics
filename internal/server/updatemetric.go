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

	name := chi.URLParam(r, "name")
	mType := metric.Type(chi.URLParam(r, "type"))
	value := chi.URLParam(r, "value")

	switch mType {
	case metric.GaugeType:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("cannot parse gauge metric: %s", err), http.StatusBadRequest)
		}
		s.Storage.SetGauge(name, v)
	case metric.CounterType:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("cannot parse counter metric: %s", err), http.StatusBadRequest)
		}
		s.Storage.AddCounter(name, v)
	default:
		http.Error(w, "unknown metric type", http.StatusNotImplemented)
		return
	}

	w.WriteHeader(http.StatusOK)
}
