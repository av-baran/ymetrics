package server

import (
	"fmt"
	"net/http"

	"github.com/av-baran/ymetrics/internal/metric"
	"github.com/av-baran/ymetrics/pkg/interrors"
	"github.com/go-chi/chi/v5"
)

func (s *Server) GetMetricHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	m := &metric.Metrics{
		ID:    chi.URLParam(r, "name"),
		MType: chi.URLParam(r, "type"),
	}
	if err := s.Storage.GetMetric(m); err != nil {
		sendError(w, "cannot get gauge metric", err)
		return
	}

	var resp string
	switch m.MType {
	case metric.GaugeType:
		resp = fmt.Sprintf("%v", *m.Value)
	case metric.CounterType:
		resp = fmt.Sprintf("%v", *m.Delta)
	default:
		sendError(w, "cannot handle get request", interrors.ErrInvalidMetricType)
		return
	}
	w.Write([]byte(resp))
}
