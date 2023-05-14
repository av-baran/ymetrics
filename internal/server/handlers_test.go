package server

import (
	"net/http/httptest"
	"testing"

	"github.com/av-baran/ymetrics/internal/repository/memstor"
	"github.com/stretchr/testify/assert"
)

func TestMetricsHandlersWithParameters(t *testing.T) {
	repo := memstor.New()
	serv := New(repo, defaultCfg)
	ts := httptest.NewServer(serv.Router)
	defer ts.Close()

	tests := handlersTestCases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, got := testRequest(t, ts, tt.method, tt.request, tt.body)
			defer resp.Body.Close()
			assert.Equal(t, tt.expectedCode, resp.StatusCode)
			if !tt.wantErr {
				assert.Equal(t, tt.expectedBody, got)
			}
		})
	}

}

func TestMetricsHandlersWithParametersAndCompression(t *testing.T) {
	repo := memstor.New()
	serv := New(repo, defaultCfg)
	ts := httptest.NewServer(serv.Router)
	defer ts.Close()

	tests := handlersTestCases
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
