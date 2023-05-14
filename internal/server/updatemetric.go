package server

import (
	"net/http"
	"strconv"

	"github.com/av-baran/ymetrics/internal/metric"
	"github.com/av-baran/ymetrics/pkg/interrors"
	"github.com/go-chi/chi/v5"
)

func (s *Server) UpdateMetricHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	m := &metric.Metric{
		ID:    chi.URLParam(r, "name"),
		MType: chi.URLParam(r, "type"),
	}
	value := chi.URLParam(r, "value")

	switch m.MType {
	case metric.GaugeType:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			sendError(w, "cannot parse gauge metric", err)
		}
		m.Value = &v
	case metric.CounterType:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			sendError(w, "cannot parse counter metric", err)
		}
		m.Delta = &v
	default:
		sendError(w, "cannot handle update request", interrors.ErrInvalidMetricType)
		return
	}
	s.Storage.SetMetric(*m)
}
