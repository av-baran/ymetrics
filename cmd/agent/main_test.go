package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/av-baran/ymetrics/internal/entities/metric"
	"github.com/stretchr/testify/assert"
)

func Test_collectMetrics(t *testing.T) {
	type testMetric struct {
		Name string
		Type metric.Type
	}
	test := struct {
		name      string
		pollCount uint64
		want      map[string]metric.Type
	}{
		name:      "Test that all necessary metrics have been received",
		pollCount: 0,
		want: map[string]metric.Type{
			"Alloc":         metric.Gauge,
			"BuckHashSys":   metric.Gauge,
			"Frees":         metric.Gauge,
			"GCCPUFraction": metric.Gauge,
			"GCSys":         metric.Gauge,
			"HeapAlloc":     metric.Gauge,
			"HeapIdle":      metric.Gauge,
			"HeapInuse":     metric.Gauge,
			"HeapObjects":   metric.Gauge,
			"HeapReleased":  metric.Gauge,
			"HeapSys":       metric.Gauge,
			"LastGC":        metric.Gauge,
			"Lookups":       metric.Gauge,
			"MCacheInuse":   metric.Gauge,
			"MCacheSys":     metric.Gauge,
			"MSpanInuse":    metric.Gauge,
			"MSpanSys":      metric.Gauge,
			"Mallocs":       metric.Gauge,
			"NextGC":        metric.Gauge,
			"NumForcedGC":   metric.Gauge,
			"NumGC":         metric.Gauge,
			"OtherSys":      metric.Gauge,
			"PauseTotalNs":  metric.Gauge,
			"StackInuse":    metric.Gauge,
			"StackSys":      metric.Gauge,
			"Sys":           metric.Gauge,
			"TotalAlloc":    metric.Gauge,
			"PollCount":     metric.Counter,
			"RandomValue":   metric.Gauge,
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
				Type:  metric.Gauge,
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
