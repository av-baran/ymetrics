package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/av-baran/ymetrics/internal/repository/memstor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	dumpTestCfg = &ServerConfig{
		ServerAddress:   "localhost:8080",
		LogLevel:        "debug",
		StoreInterval:   1,
		FileStoragePath: "/tmp/metrics-test.json",
		Restore:         true,
	}
)

func TestDump(t *testing.T) {
	firstRepo := memstor.New()
	firstServ := New(firstRepo, dumpTestCfg)
	firstTs := httptest.NewServer(firstServ.Router)
	defer firstTs.Close()

	secondRepo := memstor.New()
	secondServ := New(secondRepo, dumpTestCfg)
	secondTs := httptest.NewServer(secondServ.Router)
	defer secondTs.Close()

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
		testRequest(t, firstTs, req.method, req.request, req.body)
	}
	go firstServ.Syncfile()
	time.Sleep(time.Second * 2)

	err = secondServ.Restore()
	assert.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, got := testRequest(t, secondTs, tt.method, tt.request, tt.body)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Equal(t, tt.expectedBody, got)
		})
	}
}
