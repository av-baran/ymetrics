package server

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/av-baran/ymetrics/internal/templates"
)

func (s *Server) GetAllMetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	metrics := s.Storage.GetAllMetrics()
	strMetrics := make(map[string]string, len(metrics))

	for _, v := range metrics {
		switch v.MType {
		case "counter":
			strMetrics[v.ID] = fmt.Sprintf("%v", *v.Delta)
		case "gauge":
			strMetrics[v.ID] = fmt.Sprintf("%v", *v.Value)
		}
	}

	if err := writeTemplate(w, strMetrics); err != nil {
		http.Error(w, fmt.Sprintf("can't render metrics page: %s", err), getErrorCode(err))
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
