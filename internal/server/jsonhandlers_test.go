package server

import (
	"net/http/httptest"
	"testing"

	"github.com/av-baran/ymetrics/internal/config"
	"github.com/av-baran/ymetrics/internal/repository/memstor"
	"github.com/stretchr/testify/assert"
)

var (
	defaultCfg = &config.ServerConfig{
		ServerAddress:   "localhost:8080",
		StoreInterval:   300,
		FileStoragePath: "/tmp/metrics-db.json",
		Restore:         true,
		LoggerConfig: config.LoggerConfig{
			Level: "debug",
		},
	}
)

func TestJsonHandlers(t *testing.T) {
	repo := memstor.New()
	serv := New(repo, defaultCfg)
	ts := httptest.NewServer(serv.Router)
	defer ts.Close()

	tests := jsonHandlersTestCases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zResp, zGot := testGzRequest(t, ts, tt.method, tt.request, tt.body)
			defer zResp.Body.Close()
			assert.Equal(t, tt.expectedCode, zResp.StatusCode)
			if !tt.wantErr {
				assert.Equal(t, tt.expectedZBody, zGot)
			}
		})
	}

}

func TestJsonHandlersWithGzCompression(t *testing.T) {
	repo := memstor.New()
	serv := New(repo, defaultCfg)
	ts := httptest.NewServer(serv.Router)
	defer ts.Close()

	tests := jsonHandlersTestCases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zResp, zGot := testGzRequest(t, ts, tt.method, tt.request, tt.body)
			defer zResp.Body.Close()
			assert.Equal(t, tt.expectedCode, zResp.StatusCode)
			if !tt.wantErr {
				assert.Equal(t, tt.expectedZBody, zGot)
			}
		})
	}

}
