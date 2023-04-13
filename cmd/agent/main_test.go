package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/av-baran/ymetrics/internal/entity/metric"
	"github.com/stretchr/testify/assert"
)

func Test_collectMetrics(t *testing.T) {
	test := struct {
		pollCount uint64
		want      map[string]metric.Type
	}{
		pollCount: 0,
		want: map[string]metric.Type{
			"Alloc":         metric.GaugeType,
			"BuckHashSys":   metric.GaugeType,
			"Frees":         metric.GaugeType,
			"GCCPUFraction": metric.GaugeType,
			"GCSys":         metric.GaugeType,
			"HeapAlloc":     metric.GaugeType,
			"HeapIdle":      metric.GaugeType,
			"HeapInuse":     metric.GaugeType,
			"HeapObjects":   metric.GaugeType,
			"HeapReleased":  metric.GaugeType,
			"HeapSys":       metric.GaugeType,
			"LastGC":        metric.GaugeType,
			"Lookups":       metric.GaugeType,
			"MCacheInuse":   metric.GaugeType,
			"MCacheSys":     metric.GaugeType,
			"MSpanInuse":    metric.GaugeType,
			"MSpanSys":      metric.GaugeType,
			"Mallocs":       metric.GaugeType,
			"NextGC":        metric.GaugeType,
			"NumForcedGC":   metric.GaugeType,
			"NumGC":         metric.GaugeType,
			"OtherSys":      metric.GaugeType,
			"PauseTotalNs":  metric.GaugeType,
			"StackInuse":    metric.GaugeType,
			"StackSys":      metric.GaugeType,
			"Sys":           metric.GaugeType,
			"TotalAlloc":    metric.GaugeType,
			"PollCount":     metric.CounterType,
			"RandomValue":   metric.GaugeType,
		},
	}
	wantMetricNames := make([]string, len(test.want))
	i := 0
	for k := range test.want {
		wantMetricNames[i] = k
		i++
	}

	got := collectMetrics(test.pollCount)
	gotMetricNames := make([]string, len(got))
	for i, v := range got {
		gotMetricNames[i] = v.Name
		assert.NotNil(t, v.Value)
		assert.Equal(t, test.want[v.Name], v.Type)
	}
	assert.ElementsMatch(t, gotMetricNames, wantMetricNames)
}

func Test_sendMetric(t *testing.T) {
	tests := []struct {
		name     string
		inMetric inMetric
		wantErr  bool
	}{
		{
			name: "Wrong type",
			inMetric: inMetric{
				Name:  "unknownMetric",
				Value: 22,
				Type:  metric.Type("unknown"),
			},
			wantErr: true,
		},
		{
			name: "Correct type",
			inMetric: inMetric{
				Name:  "unknownMetric",
				Value: 22,
				Type:  metric.GaugeType,
			},
			wantErr: false,
		},
	}
	srv := httpMock("/update/", http.StatusOK, "Ok")
	defer srv.Close()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := sendMetric(srv.URL, tt.inMetric); (err != nil) != tt.wantErr {
				t.Errorf("sendMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func httpMock(path string, statusCode int, resp interface{}) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc(path, dummyHandler(statusCode, resp))
	return httptest.NewServer(mux)
}

func dummyHandler(statusCode int, resp interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		fmt.Fprint(w, resp)
	}
}
