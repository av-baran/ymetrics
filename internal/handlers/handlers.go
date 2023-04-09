package handlers

import (
	"log"
	"net/http"

	"github.com/av-baran/ymetrics/internal/entities/metric"
	"github.com/av-baran/ymetrics/internal/interrors"
	"github.com/go-chi/chi/v5"
)

type Storage interface {
	UpdateMetric(*metric.Rawdata) error
	GetMetric(*metric.Rawdata) (string, error)
}

func UpdateMetricHandler(s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		var m = &metric.Rawdata{
			Name:  chi.URLParam(r, "name"),
			Type:  metric.Type(chi.URLParam(r, "type")),
			Value: chi.URLParam(r, "value"),
		}

		if err := s.UpdateMetric(m); err != nil {
			log.Printf(`Error "%v" while updating metric`, err.Error())
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

func GetMetricHandler(s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		var m = &metric.Rawdata{
			Name: chi.URLParam(r, "name"),
			Type: metric.Type(chi.URLParam(r, "type")),
		}

		mValue, err := s.GetMetric(m)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		w.Write([]byte(mValue))
	}
}
