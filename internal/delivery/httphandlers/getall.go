package httphandlers

import (
	"fmt"
	"net/http"
)

const (
	pageHeader = `<html><body><h1>Metric list</h1><ul>`
	pageFooter = `</ul></body></html>`
)

func GetAllMetricsHandler(g metricGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		allMetrics := g.GetAllMetrics()
		fmt.Fprint(w, pageHeader)
		for k, v := range allMetrics {
			fmt.Fprintf(w, "<li> %v = %v </li>", k, v)
		}
		fmt.Fprint(w, pageFooter)
	}
}
