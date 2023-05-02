package agent

// import (
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"net/url"
// 	"testing"
//
// 	"github.com/av-baran/ymetrics/internal/config"
// 	"github.com/av-baran/ymetrics/internal/metric"
// 	"github.com/av-baran/ymetrics/pkg/interrors"
// 	"github.com/stretchr/testify/assert"
// )
//
// func Test_collectMetrics(t *testing.T) {
// 	test := struct {
// 		want map[string]metric.Type
// 	}{
// 		want: map[string]metric.Type{
// 			"Alloc":         metric.GaugeType,
// 			"BuckHashSys":   metric.GaugeType,
// 			"Frees":         metric.GaugeType,
// 			"GCCPUFraction": metric.GaugeType,
// 			"GCSys":         metric.GaugeType,
// 			"HeapAlloc":     metric.GaugeType,
// 			"HeapIdle":      metric.GaugeType,
// 			"HeapInuse":     metric.GaugeType,
// 			"HeapObjects":   metric.GaugeType,
// 			"HeapReleased":  metric.GaugeType,
// 			"HeapSys":       metric.GaugeType,
// 			"LastGC":        metric.GaugeType,
// 			"Lookups":       metric.GaugeType,
// 			"MCacheInuse":   metric.GaugeType,
// 			"MCacheSys":     metric.GaugeType,
// 			"MSpanInuse":    metric.GaugeType,
// 			"MSpanSys":      metric.GaugeType,
// 			"Mallocs":       metric.GaugeType,
// 			"NextGC":        metric.GaugeType,
// 			"NumForcedGC":   metric.GaugeType,
// 			"NumGC":         metric.GaugeType,
// 			"OtherSys":      metric.GaugeType,
// 			"PauseTotalNs":  metric.GaugeType,
// 			"StackInuse":    metric.GaugeType,
// 			"StackSys":      metric.GaugeType,
// 			"Sys":           metric.GaugeType,
// 			"TotalAlloc":    metric.GaugeType,
// 			"PollCount":     metric.CounterType,
// 			"RandomValue":   metric.GaugeType,
// 		},
// 	}
// 	wantMetricNames := make([]string, len(test.want))
// 	i := 0
// 	for k := range test.want {
// 		wantMetricNames[i] = k
// 		i++
// 	}
// 	cfg := config.NewAgentConfig()
// 	a := NewAgent(cfg)
//
// 	a.collectMetrics()
// 	gotMetricNames := make([]string, len(collectedMetrics))
// 	for i, v := range collectedMetrics {
// 		gotMetricNames[i] = v.Name
// 		assert.NotNil(t, v.Value)
// 		assert.Equal(t, test.want[v.Name], v.Type)
// 	}
// 	assert.ElementsMatch(t, gotMetricNames, wantMetricNames)
// }
//
// func Test_sendMetricOk(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		Metric  metric.Metric
// 		wantErr bool
// 	}{
// 		{
// 			name: "Wrong type",
// 			Metric: metric.Metric{
// 				Name:  "unknownMetric",
// 				Value: 22,
// 				Type:  metric.Type("unknown"),
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "Correct type",
// 			Metric: metric.Metric{
// 				Name:  "unknownMetric",
// 				Value: 22,
// 				Type:  metric.GaugeType,
// 			},
// 			wantErr: false,
// 		},
// 	}
//
// 	srv := httpServerMock("/update/", http.StatusOK, "Ok")
// 	defer srv.Close()
// 	u, _ := url.Parse(srv.URL)
//
// 	cfg := &config.AgentConfig{
// 		ServerAddress:  fmt.Sprintf("%v:%v", u.Hostname(), u.Port()),
// 		PollInterval:   1,
// 		ReportInterval: 2,
// 	}
//
// 	a := NewAgent(cfg)
// 	a.collectMetrics()
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got1 := a.sendMetric(tt.Metric)
// 			if tt.wantErr {
// 				assert.Error(t, got1)
// 			} else {
// 				assert.NoError(t, got1)
// 			}
// 		})
// 	}
// }
//
// func Test_sendMetricErr(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		Metric  metric.Metric
// 		wantErr bool
// 	}{
// 		{
// 			name: "Wrong type",
// 			Metric: metric.Metric{
// 				Name:  "unknownMetric",
// 				Value: 22,
// 				Type:  metric.Type("unknown"),
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "Correct type",
// 			Metric: metric.Metric{
// 				Name:  "unknownMetric",
// 				Value: 22,
// 				Type:  metric.GaugeType,
// 			},
// 			wantErr: false,
// 		},
// 	}
//
// 	srv := httpServerMock("/update/", http.StatusInternalServerError, interrors.ErrStorageInternalError.Error()+"\n")
// 	defer srv.Close()
// 	u, _ := url.Parse(srv.URL)
//
// 	cfg := &config.AgentConfig{
// 		ServerAddress:  fmt.Sprintf("%v:%v", u.Hostname(), u.Port()),
// 		PollInterval:   1,
// 		ReportInterval: 2,
// 	}
//
// 	a := NewAgent(cfg)
// 	a.collectMetrics()
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := a.sendMetric(tt.Metric)
// 			assert.Error(t, got)
// 		})
// 	}
// }
//
// func TestRun(t *testing.T) {
// 	srv := httpServerMock("/update/", http.StatusOK, "Ok")
// 	defer srv.Close()
// 	u, _ := url.Parse(srv.URL)
//
// 	cfg := &config.AgentConfig{
// 		ServerAddress:  fmt.Sprintf("%v:%v", u.Hostname(), u.Port()),
// 		PollInterval:   1,
// 		ReportInterval: 2,
// 	}
//
// 	okAgent := NewAgent(cfg)
// 	go okAgent.Run(cfg)
// }
//
// func Test_dumpOk(t *testing.T) {
// 	srv := httpServerMock("/update/", http.StatusOK, "Ok")
// 	defer srv.Close()
// 	u, _ := url.Parse(srv.URL)
//
// 	cfg := &config.AgentConfig{
// 		ServerAddress:  fmt.Sprintf("%v:%v", u.Hostname(), u.Port()),
// 		PollInterval:   1,
// 		ReportInterval: 2,
// 	}
//
// 	a := NewAgent(cfg)
// 	a.collectMetrics()
// 	err := a.dump()
// 	assert.NoError(t, err)
// }
//
// func Test_dumpErr(t *testing.T) {
// 	srv := httpServerMock("/update/", http.StatusInternalServerError, interrors.ErrStorageInternalError.Error()+"\n")
// 	defer srv.Close()
// 	u, _ := url.Parse(srv.URL)
//
// 	cfg := &config.AgentConfig{
// 		ServerAddress:  fmt.Sprintf("%v:%v", u.Hostname(), u.Port()),
// 		PollInterval:   1,
// 		ReportInterval: 2,
// 	}
//
// 	a := NewAgent(cfg)
// 	a.collectMetrics()
// 	err := a.dump()
// 	assert.Error(t, err)
// }
//
// func httpServerMock(path string, statusCode int, resp interface{}) *httptest.Server {
// 	mux := http.NewServeMux()
// 	mux.HandleFunc(path, dummyHandler(statusCode, resp))
// 	return httptest.NewServer(mux)
// }
//
// func dummyHandler(statusCode int, resp interface{}) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(statusCode)
// 		fmt.Fprint(w, resp)
// 	}
// }
