package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/av-baran/ymetrics/internal/entities/metric"
	"github.com/av-baran/ymetrics/internal/interrors"
)

type storage interface {
	UpdateMetric(*metric.Rawdata) error
}

func UpdateMetricHandler(s storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request with url: %v", r.URL.String())
		if r.Method != http.MethodPost {
			log.Printf("Method %v is not allowed.", r.Method)
			http.Error(w, "Only POST request is allowed.", http.StatusMethodNotAllowed)
			return
		}
		metric, err := parseURL(r.URL.Path)
		if err != nil {
			log.Printf("Error while parsing url")
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err := s.UpdateMetric(metric); err != nil {
			log.Printf("Error while updating metric")
			var statusCode int
			switch err.Error() {
			case interrors.ErrInvalidMetricType:
				statusCode = http.StatusNotImplemented
			case interrors.ErrInvalidMetricValue:
				statusCode = http.StatusBadRequest
			case interrors.ErrMetricAlreadyExists:
				statusCode = http.StatusBadRequest
			default:
				statusCode = http.StatusInternalServerError
			}
			http.Error(w, err.Error(), statusCode)
			return
		}
	}
}

func parseURL(path string) (*metric.Rawdata, error) {
	p := strings.Split(path, "/")
	if len(p) != 5 {
		return nil, errors.New(interrors.ErrBadURL)
	}
	return &metric.Rawdata{
		Type:  metric.Type(p[2]),
		Name:  p[3],
		Value: p[4],
	}, nil
}
