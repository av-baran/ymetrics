package server

import (
	"net/http"

	"github.com/av-baran/ymetrics/internal/metric"
	"github.com/go-chi/chi/v5"
)

func (s *Server) GetMetricHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	name := chi.URLParam(r, "name")
	mType := metric.Type(chi.URLParam(r, "type"))
	value := ""

	switch mType {
	case metric.GaugeType:
		value = s.Storage.GetGauge(name)
	case metric.CounterType:
		value = s.Storage.GetCounter(name)
	}
	w.Write([]byte(value))
}
