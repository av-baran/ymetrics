package server

import (
	"fmt"
	"net/http"
)

//text/template
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

	writeHTML(w, metrics)
}

func writeHTML(w http.ResponseWriter, metrics map[string]string) {
	const (
		pageHeader = `<html><body><h1>Metric list</h1><ul>`
		pageFooter = `</ul></body></html>`
	)

	fmt.Fprint(w, pageHeader)

	for k, v := range metrics {
		fmt.Fprintf(w, "<li> %v = %v </li>", k, v)
	}

	fmt.Fprint(w, pageFooter)
}
