package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/av-baran/ymetrics/internal/metric"
)

func (s *Server) UpdateBatchJSONHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.HandlerTimeout)
	defer cancel()

	readBody, err := io.ReadAll(r.Body)
	if err != nil {
		sendError(w, "cannot read request body", err)
		return
	}
	r.Body.Close()

	metrics := make([]metric.Metric, 0)

	if err := json.Unmarshal(readBody, &metrics); err != nil {
		sendError(w, "cannot unmarshal request body", err)
		return
	}

	if err := s.Storage.SetMetricsBatch(ctx, metrics); err != nil {
		sendError(w, "cannot set metrics batch", err)
		return
	}
}
