package httphandlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/av-baran/ymetrics/internal/entity/metric"
	"github.com/av-baran/ymetrics/pkg/interrors"
	"github.com/go-chi/chi/v5"
)

func UpdateMetricHandler(u metricUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		var m = &metric.Rawdata{
			Name:  chi.URLParam(r, "name"),
			Type:  metric.Type(chi.URLParam(r, "type")),
			Value: chi.URLParam(r, "value"),
		}

		if err := storeMetric(u, m); err != nil {
			fmt.Printf(`Error "%v" while updating metric`, err.Error())
			statusCode := getErrorCode(err)
			http.Error(w, err.Error(), statusCode)
		}
	}
}

func storeMetric(u metricUpdater, m *metric.Rawdata) error {
	switch m.Type {
	case metric.GaugeType:
		v, err := strconv.ParseFloat(m.Value, 64)
		if err != nil {
			return errors.New(interrors.ErrInvalidMetricValue)
		}
		if err := u.UpdateGauge(&metric.Gauge{Name: m.Name, Value: v}); err != nil {
			return err
		}
		return nil
	case metric.CounterType:
		v, err := strconv.ParseInt(m.Value, 10, 64)
		if err != nil {
			return errors.New(interrors.ErrInvalidMetricValue)
		}
		if err := u.UpdateCounter(&metric.Counter{Name: m.Name, Value: v}); err != nil {
			return err
		}
		return nil
	default:
		return errors.New(interrors.ErrInvalidMetricType)
	}
}
