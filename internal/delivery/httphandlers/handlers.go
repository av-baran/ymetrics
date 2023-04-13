package httphandlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/av-baran/ymetrics/internal/entity/metric"
	"github.com/av-baran/ymetrics/pkg/interrors"
	"github.com/go-chi/chi/v5"
)

type metricUpdater interface {
	UpdateGauge(*metric.Gauge) error
	UpdateCounter(*metric.Counter) error
}

type metricGetter interface {
	GetGauge(string) (float64, error)
	GetCounter(string) (int64, error)
}

func UpdateMetricHandler(u metricUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		var m = &metric.Rawdata{
			Name:  chi.URLParam(r, "name"),
			Type:  metric.Type(chi.URLParam(r, "type")),
			Value: chi.URLParam(r, "value"),
		}

		if err := storeMetric(u, m); err != nil {
			log.Printf(`Error "%v" while updating metric`, err.Error())
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

func getErrorCode(e error) (statusCode int) {
	switch e.Error() {
	case interrors.ErrInvalidMetricType:
		statusCode = http.StatusNotImplemented
	case interrors.ErrInvalidMetricValue:
		statusCode = http.StatusBadRequest
	case interrors.ErrMetricExistsWithAnotherType:
		statusCode = http.StatusBadRequest
	case interrors.ErrMetricNotFound:
		statusCode = http.StatusNotFound
	default:
		statusCode = http.StatusInternalServerError
	}
	return
}

func GetAllMetricsHandler(g metricGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
