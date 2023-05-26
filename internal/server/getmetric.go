package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/av-baran/ymetrics/internal/logger"
	"github.com/av-baran/ymetrics/internal/metric"
	"github.com/av-baran/ymetrics/pkg/interrors"
	"github.com/go-chi/chi/v5"
)

func (s *Server) GetMetricHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.HandlerTimeout)
	defer cancel()

	w.Header().Set("Content-Type", "text/plain")

	metricID := chi.URLParam(r, "name")
	metricType := chi.URLParam(r, "type")
	m, err := s.Storage.GetMetric(ctx, metricID, metricType)
	if err != nil {
		sendError(w, "cannot get metric", err)
		return
	}

	logger.Errorf("get metric handler m: %+v", m)
	var resp string
	switch metricType {
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
