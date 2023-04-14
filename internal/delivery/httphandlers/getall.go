package httphandlers

import (
	"net/http"
)

func GetAllMetricsHandler(g metricGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
