package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/av-baran/ymetrics/internal/metric"
	"github.com/av-baran/ymetrics/pkg/interrors"
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
			http.Error(w, fmt.Sprintf("cannot parse gauge metric: %s", err), getErrorCode(err))
		}
		m.Value = &v
	case metric.CounterType:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("cannot parse counter metric: %s", err), getErrorCode(err))
		}
		m.Delta = &v
	default:
		err := interrors.ErrInvalidMetricType
		http.Error(w, err.Error(), getErrorCode(err))
		return
	}
	s.Storage.SetMetric(*m)
}
