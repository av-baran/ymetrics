package server

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/av-baran/ymetrics/internal/templates"
)

func (s *Server) GetAllMetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	counters := s.Storage.GetAllCounter()
	gauges := s.Storage.GetAllGauge()

	metrics := make(map[string]string, len(counters)+len(gauges))

	for k, v := range counters {
		metrics[k] = fmt.Sprintf("%v", v)
	}
	for k, v := range gauges {
		metrics[k] = fmt.Sprintf("%v", v)
	}

	if err := writeTemplate(w, metrics); err != nil {
		http.Error(w, fmt.Sprintf("can't render metrics page: %s", err), http.StatusInternalServerError)
	}
}

func writeTemplate(w http.ResponseWriter, metrics map[string]string) error {
	t := template.New("t")

	t, err := t.Parse(templates.GetAllPageTemplate)
	if err != nil {
		return fmt.Errorf("cannot parse template: %w", err)
	}

	if err := t.Execute(w, metrics); err != nil {
		return fmt.Errorf("cannot render template: %w", err)
	}

	return nil
}
