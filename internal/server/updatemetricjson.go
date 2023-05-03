package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/av-baran/ymetrics/internal/metric"
)

func (s *Server) UpdateMetricJSONHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	readBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot read request body: %s", err), http.StatusBadRequest)
		return
	}
	r.Body.Close()

	m := &metric.Metrics{}
	if err := json.Unmarshal(readBody, m); err != nil {
		http.Error(w, fmt.Sprintf("cannot unmarshal request body: %s", err), http.StatusBadRequest)
		return
	}

	if err := s.Storage.SetMetric(*m); err != nil {
		http.Error(w, fmt.Sprintf("cannot set metric: %s", err), http.StatusNotFound)
		return
	}

	resM := &metric.Metrics{
		ID:    m.ID,
		MType: m.MType,
	}
	if err := s.Storage.GetMetric(resM); err != nil {
		http.Error(w, fmt.Sprintf("cannot get metric: %s", err), http.StatusNotFound)
		return
	}

	respBody, err := json.Marshal(&resM)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot marshal response body: %s", err), http.StatusInternalServerError)
		return
	}

	w.Write(respBody)
	r.Body = io.NopCloser(bytes.NewReader(readBody))
}
