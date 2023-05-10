package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/av-baran/ymetrics/internal/config"
	"github.com/av-baran/ymetrics/internal/repository/memstor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	dumpTestCfg = &config.ServerConfig{
		ServerAddress:   "localhost:8080",
		StoreInterval:   1,
		FileStoragePath: "/tmp/metrics-test.json",
		Restore:         true,
		LoggerConfig: config.LoggerConfig{
			Level: "debug",
		},
	}
)

func TestDump(t *testing.T) {
	firstRepo := memstor.New()
	firstServ := New(firstRepo, dumpTestCfg)
	firstTS := httptest.NewServer(firstServ.Router)
	defer firstTS.Close()
	defer firstServ.Shutdown()

	secondRepo := memstor.New()
	secondServ := New(secondRepo, dumpTestCfg)
	secondTS := httptest.NewServer(secondServ.Router)
	defer secondTS.Close()
	defer secondServ.Shutdown()

	f, err := os.Create(dumpTestCfg.FileStoragePath)
	require.NoError(t, err)
	f.Close()
	defer os.Remove(dumpTestCfg.FileStoragePath)

	data := []struct {
		request string
		method  string
		body    string
	}{
		{
			request: "/update/",
			method:  http.MethodPost,
			body:    `{"id":"someCounter","type":"counter","delta":5}`,
		},
		{
			request: "/update/",
			method:  http.MethodPost,
			body:    `{"id":"someGauge","type":"gauge","value":0.1}`,
		},
		{
			request: "/update/",
			method:  http.MethodPost,
			body:    `{"id":"anotherGauge","type":"gauge","value":5555}`,
		},
	}

	tests := []struct {
		name         string
		request      string
		method       string
		body         string
		expectedBody string
	}{
		{
			name:         "get counter",
			request:      "/value/",
			method:       http.MethodPost,
			body:         `{"id":"someCounter","type":"counter"}`,
			expectedBody: `{"id":"someCounter","type":"counter","delta":5}`,
		},
		{
			name:         "get gauge",
			request:      "/value/",
			method:       http.MethodPost,
			body:         `{"id":"someGauge","type":"gauge"}`,
			expectedBody: `{"id":"someGauge","type":"gauge","value":0.1}`,
		},
		{
			name:         "get another gauge",
			request:      "/value/",
			method:       http.MethodPost,
			body:         `{"id":"anotherGauge","type":"gauge"}`,
			expectedBody: `{"id":"anotherGauge","type":"gauge","value":5555}`,
		},
	}

	for _, req := range data {
		resp, _ := testRequest(t, firstTS, req.method, req.request, req.body)
		resp.Body.Close()
	}
	go firstServ.syncfile()
	time.Sleep(time.Second * 2)

	err = secondServ.restore()
	assert.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, got := testRequest(t, secondTS, tt.method, tt.request, tt.body)
			resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Equal(t, tt.expectedBody, got)
		})
	}
}
