package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/av-baran/ymetrics/internal/metric"
)

func (s *Server) GetMetricJSONHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	readBody, err := io.ReadAll(r.Body)
	if err != nil {
		sendError(w, "cannot read request body", err)
		return
	}
	r.Body.Close()

	m := &metric.Metric{}
	if err := json.Unmarshal(readBody, m); err != nil {
		sendError(w, "cannot unmarshal request body", err)
		return
	}

	respM, err := s.Storage.GetMetric(m.ID, m.MType)
	if err != nil {
		sendError(w, "cannot get metric", err)
		return
	}

	respBody, err := json.Marshal(respM)
	if err != nil {
		sendError(w, "cannot marshal response body", err)
		return
	}

	w.Write(respBody)

	r.Body = io.NopCloser(bytes.NewReader(readBody))
}
