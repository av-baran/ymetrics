package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/av-baran/ymetrics/internal/metric"
)

func (s *Server) UpdateMetricJSONHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	readBody, err := io.ReadAll(r.Body)
	if err != nil {
		sendError(w, "cannot read request body", err)
		return
	}
	r.Body.Close()

	m := &metric.Metrics{}
	if err := json.Unmarshal(readBody, m); err != nil {
		sendError(w, "cannot unmarshal request body", err)
		return
	}

	if err := s.Storage.SetMetric(*m); err != nil {
		sendError(w, "cannot set metric", err)
		return
	}

	resM := &metric.Metrics{
		ID:    m.ID,
		MType: m.MType,
	}
	if err := s.Storage.GetMetric(resM); err != nil {
		sendError(w, "cannot get metric", err)
		return
	}

	respBody, err := json.Marshal(&resM)
	if err != nil {
		sendError(w, "cannot marshal response body", err)
		return
	}

	w.Write(respBody)
	r.Body = io.NopCloser(bytes.NewReader(readBody))
}
