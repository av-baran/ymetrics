package agent

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/av-baran/ymetrics/internal/config"
	"github.com/av-baran/ymetrics/internal/metric"
	"github.com/av-baran/ymetrics/pkg/interrors"
	"github.com/stretchr/testify/assert"
)

var (
	defaultCfg = &config.AgentConfig{
		ServerAddress:  "localhost:8080",
		PollInterval:   3,
		ReportInterval: 5,
		RetryConfig: config.RetryConfig{
			BackoffLimit:  1,
			SleepTime:     1,
			TimeIncrement: 1,
		},
	}

	testCfg = config.AgentConfig{
		PollInterval:   1,
		ReportInterval: 2,
		RetryConfig: config.RetryConfig{
			BackoffLimit:  1,
			SleepTime:     1,
			TimeIncrement: 1,
		},
	}
)

func Test_collectMetrics(t *testing.T) {
	test := struct {
		want map[string]string
	}{
		want: map[string]string{
			"Alloc":           metric.GaugeType,
			"BuckHashSys":     metric.GaugeType,
			"Frees":           metric.GaugeType,
			"GCCPUFraction":   metric.GaugeType,
			"GCSys":           metric.GaugeType,
			"HeapAlloc":       metric.GaugeType,
			"HeapIdle":        metric.GaugeType,
			"HeapInuse":       metric.GaugeType,
			"HeapObjects":     metric.GaugeType,
			"HeapReleased":    metric.GaugeType,
			"HeapSys":         metric.GaugeType,
			"LastGC":          metric.GaugeType,
			"Lookups":         metric.GaugeType,
			"MCacheInuse":     metric.GaugeType,
			"MCacheSys":       metric.GaugeType,
			"MSpanInuse":      metric.GaugeType,
			"MSpanSys":        metric.GaugeType,
			"Mallocs":         metric.GaugeType,
			"NextGC":          metric.GaugeType,
			"NumForcedGC":     metric.GaugeType,
			"NumGC":           metric.GaugeType,
			"OtherSys":        metric.GaugeType,
			"PauseTotalNs":    metric.GaugeType,
			"StackInuse":      metric.GaugeType,
			"StackSys":        metric.GaugeType,
			"Sys":             metric.GaugeType,
			"TotalAlloc":      metric.GaugeType,
			"PollCount":       metric.CounterType,
			"RandomValue":     metric.GaugeType,
			"TotalMemory":     metric.GaugeType,
			"FreeMemory":      metric.GaugeType,
			"CPUutilization1": metric.GaugeType,
		},
	}
	wantMetricNames := make([]string, len(test.want))
	i := 0
	for k := range test.want {
		wantMetricNames[i] = k
		i++
	}
	a := NewAgent(defaultCfg)

	a.collectMetrics()
	gotMetricNames := make([]string, len(collectedMetrics))
	for i, v := range collectedMetrics {
		gotMetricNames[i] = v.ID
		switch v.MType {
		case metric.GaugeType:
			assert.NotNil(t, v.Value)
		case metric.CounterType:
			assert.NotNil(t, v.Delta)
			assert.Equal(t, test.want[v.ID], v.MType)
		}
	}
	assert.ElementsMatch(t, gotMetricNames, wantMetricNames)
}

func Test_sendMetricOk(t *testing.T) {
	tests := []struct {
		name    string
		Metric  []metric.Metric
		wantErr bool
	}{
		{
			name: "Correct type",
			Metric: []metric.Metric{
				{
					ID:    "knownMetric",
					Value: getFloat64Ptr(22),
					MType: metric.GaugeType,
				},
			},
			wantErr: false,
		},
	}

	srv := httpServerMock("/updates/", http.StatusOK, "Ok")
	defer srv.Close()
	u, _ := url.Parse(srv.URL)
	tCfg := testCfg
	tCfg.ServerAddress = fmt.Sprintf("%v:%v", u.Hostname(), u.Port())
	a := NewAgent(&tCfg)

	a.collectMetrics()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1 := a.sendBatchJSON(tt.Metric)
			if tt.wantErr {
				assert.Error(t, got1)
			} else {
				assert.NoError(t, got1)
			}
		})
	}
}

func Test_sendMetricErr(t *testing.T) {
	tests := []struct {
		name    string
		Metric  []metric.Metric
		wantErr bool
	}{
		{
			name: "Correct type",
			Metric: []metric.Metric{
				{
					ID:    "unknownMetric",
					Value: getFloat64Ptr(22),
					MType: metric.GaugeType,
				},
			},
			wantErr: false,
		},
	}

	srv := httpServerMock("/updates/", http.StatusInternalServerError, interrors.ErrStorageInternalError.Error()+"\n")
	defer srv.Close()
	u, _ := url.Parse(srv.URL)
	tCfg := testCfg
	tCfg.ServerAddress = fmt.Sprintf("%v:%v", u.Hostname(), u.Port())
	a := NewAgent(&tCfg)

	a.collectMetrics()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := a.sendBatchJSON(tt.Metric)
			assert.Error(t, got)
		})
	}
}

func TestRun(t *testing.T) {
	srv := httpServerMock("/updates/", http.StatusOK, "Ok")
	defer srv.Close()
	u, _ := url.Parse(srv.URL)
	tCfg := testCfg
	tCfg.ServerAddress = fmt.Sprintf("%v:%v", u.Hostname(), u.Port())

	okAgent := NewAgent(&tCfg)
	go okAgent.Run()
}

func Test_dumpOk(t *testing.T) {
	srv := httpServerMock("/updates/", http.StatusOK, "Ok")
	defer srv.Close()
	u, _ := url.Parse(srv.URL)
	tCfg := testCfg
	tCfg.ServerAddress = fmt.Sprintf("%v:%v", u.Hostname(), u.Port())

	a := NewAgent(&tCfg)
	a.collectMetrics()
	err := a.batchDump()
	assert.NoError(t, err)
}

func Test_dumpErr(t *testing.T) {
	srv := httpServerMock("/updates/", http.StatusInternalServerError, interrors.ErrStorageInternalError.Error()+"\n")
	defer srv.Close()
	u, _ := url.Parse(srv.URL)
	tCfg := testCfg
	tCfg.ServerAddress = fmt.Sprintf("%v:%v", u.Hostname(), u.Port())

	a := NewAgent(&tCfg)

	a.collectMetrics()
	err := a.batchDump()
	assert.Error(t, err)
}

func httpServerMock(path string, statusCode int, resp interface{}) *httptest.Server {
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
