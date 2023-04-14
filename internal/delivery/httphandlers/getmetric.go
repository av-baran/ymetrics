package httphandlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/av-baran/ymetrics/internal/entity/metric"
	"github.com/av-baran/ymetrics/pkg/interrors"
	"github.com/go-chi/chi/v5"
)

func GetMetricHandler(g metricGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		var m = &metric.Rawdata{
			Name: chi.URLParam(r, "name"),
			Type: metric.Type(chi.URLParam(r, "type")),
		}

		if mValue, err := getMetric(g, m); err != nil {
			log.Printf(`Error "%v" while updating metric`, err.Error())
			statusCode := getErrorCode(err)
			http.Error(w, err.Error(), statusCode)
		} else {
			w.Write([]byte(mValue))
		}
	}
}

func getMetric(g metricGetter, m *metric.Rawdata) (string, error) {
	switch m.Type {
	case metric.GaugeType:
		v, err := g.GetGauge(m.Name)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%v", v), nil
	case metric.CounterType:
		v, err := g.GetCounter(m.Name)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%v", v), nil
	default:
		return "", errors.New(interrors.ErrInvalidMetricType)
	}
}
