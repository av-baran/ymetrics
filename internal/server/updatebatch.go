package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/av-baran/ymetrics/internal/metric"
)

func (s *Server) UpdateBatchJSONHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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

	if err := s.Storage.UpdateBatch(metrics); err != nil {
		sendError(w, "cannot update metrics", err)
		return
	}

}
