package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/av-baran/ymetrics/internal/entities/metric"
	"github.com/av-baran/ymetrics/internal/httperror"
)

type storage interface {
	UpdateMetric(*metric.Rawdata) *httperror.Error
}

func UpdateMetrics(s storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request with url: %v", r.URL.String())
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST request is allowed.", http.StatusMethodNotAllowed)
			return
		}
		metric, err := parseURL(r.URL.Path)
		if err != nil {
			log.Printf("Error while parsing url")
			http.Error(w, err.Msg, err.Code)
			return
		}
		if err := s.UpdateMetric(metric); err != nil {
			log.Printf("Error while updating metric")
			http.Error(w, err.Msg, err.Code)
			return
		}
	}
}

func parseURL(path string) (*metric.Rawdata, *httperror.Error) {
	p := strings.Split(path, "/")
	if len(p) != 5 {
		return nil, httperror.New("bad request: malformed URL", http.StatusNotFound)
	}
	return &metric.Rawdata{
		Type:  metric.Type(p[2]),
		Name:  p[3],
		Value: p[4],
	}, nil
}
