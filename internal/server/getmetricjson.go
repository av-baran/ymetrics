package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/av-baran/ymetrics/internal/metric"
)

func (s *Server) GetMetricJSONHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	readBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot read request body: %s", err), http.StatusBadRequest)
		return
	}
	r.Body.Close()

	m := metric.Metrics{}
	if err := json.Unmarshal(readBody, &m); err != nil {
		http.Error(w, fmt.Sprintf("cannot unmarshal request body: %s", err), http.StatusBadRequest)
		return
	}

	switch m.MType {
	case "gauge":
		v, err := s.Storage.GetGauge(m.ID)
		if err != nil {
			http.Error(w, fmt.Sprintf("cannot get gauge: %s", err), http.StatusNotFound)
			return
		}
		m.Value = &v
	case "counter":
		d, err := s.Storage.GetCounter(m.ID)
		if err != nil {
			http.Error(w, fmt.Sprintf("cannot get counter: %s", err), http.StatusNotFound)
			return
		}
		m.Delta = &d
	default:
		http.Error(w, fmt.Sprintf("unknown metric type"), http.StatusNotImplemented)
		return
	}

	respBody, err := json.Marshal(&m)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot marshal response body: %s", err), http.StatusInternalServerError)
		return
	}

	w.Write(respBody)

	r.Body = io.NopCloser(bytes.NewReader(readBody))
	w.WriteHeader(http.StatusOK)
}
