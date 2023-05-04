package server

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/av-baran/ymetrics/internal/repository/memstor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	defaultCfg = &ServerConfig{
		ServerAddress:   "localhost:8080",
		LogLevel:        "debug",
		StoreInterval:   300,
		FileStoragePath: "/tmp/metrics-db.json",
		Restore:         true,
	}
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body string) (*http.Response, string) {
	b := bytes.NewBuffer([]byte(body))
	req, err := http.NewRequest(method, ts.URL+path, b)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func testGzRequest(t *testing.T, ts *httptest.Server, method, path string, body string) (*http.Response, string) {
	var zb bytes.Buffer
	gzw := gzip.NewWriter(&zb)

	_, err := gzw.Write([]byte(body))
	require.NoError(t, err)
	gzw.Close()

	req, err := http.NewRequest(method, ts.URL+path, &zb)
	require.NoError(t, err)

	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Encoding", "gzip")

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	zRespBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	gzr, err := gzip.NewReader(bytes.NewReader(zRespBody))
	require.NoError(t, err)
	defer gzr.Close()

	var respBody bytes.Buffer
	_, err = respBody.ReadFrom(gzr)
	require.NoError(t, err)

	return resp, respBody.String()
}

func TestNew(t *testing.T) {
	s := memstor.New()
	newServer := New(s, defaultCfg)
	assert.NotEmpty(t, newServer)
}

func TestServer(t *testing.T) {
	repo := memstor.New()
	serv := New(repo, defaultCfg)
	ts := httptest.NewServer(serv.Router)
	defer ts.Close()

	tests := []struct {
		name          string
		request       string
		method        string
		body          string
		expectedCode  int
		expectedBody  string
		expectedZBody string
		wantErr       bool
	}{
		{
			name:         "GET params",
			request:      "/update/gauge/name/1.0",
			body:         "",
			method:       http.MethodGet,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "",
			wantErr:      true,
		},
		{
			name:         "PUT params",
			request:      "/update/gauge/name/1.0",
			body:         "",
			method:       http.MethodPut,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "",
			wantErr:      true,
		},
		{
			name:         "DELETE params",
			request:      "/update/gauge/name/1.0",
			body:         "",
			method:       http.MethodDelete,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "",
			wantErr:      true,
		},
		{
			name:         "POST - params update bad type",
			request:      "/update/unknowntype/name/1.0",
			body:         "",
			method:       http.MethodPost,
			expectedCode: http.StatusNotImplemented,
			expectedBody: "",
			wantErr:      true,
		},
		{
			name:         "POST - params update gauge with bad value",
			request:      "/update/gauge/name/fuuuu",
			body:         "",
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			expectedBody: "",
			wantErr:      true,
		},
		{
			name:         "POST - params update counter with bad value",
			request:      "/update/counter/name/fuuuu",
			body:         "",
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			expectedBody: "",
			wantErr:      true,
		},
		{
			name:         "POST - params update gauge ok",
			request:      "/update/gauge/name/1.01",
			body:         "",
			method:       http.MethodPost,
			expectedCode: http.StatusOK,
			expectedBody: "",
			wantErr:      false,
		},
		{
			name:         "POST - params update counter ok",
			request:      "/update/counter/name/1",
			body:         "",
			method:       http.MethodPost,
			expectedCode: http.StatusOK,
			expectedBody: "",
			wantErr:      false,
		},
		{
			name:          "GET - params get gauge ok",
			request:       "/value/gauge/name",
			body:          "",
			method:        http.MethodGet,
			expectedCode:  http.StatusOK,
			expectedBody:  "1.01",
			expectedZBody: "1.01",
			wantErr:       false,
		},
		{
			name:          "GET - params get counter ok",
			request:       "/value/counter/name",
			body:          "",
			method:        http.MethodGet,
			expectedCode:  http.StatusOK,
			expectedBody:  "1",
			expectedZBody: "2",
			wantErr:       false,
		},
		{
			name:         "GET - params get unknown gauge",
			request:      "/value/gauge/unknownname",
			body:         "",
			method:       http.MethodGet,
			expectedCode: http.StatusNotFound,
			expectedBody: "",
			wantErr:      true,
		},
		{
			name:         "GET - params get unknown counter",
			request:      "/value/counter/unknownname",
			body:         "",
			method:       http.MethodGet,
			expectedCode: http.StatusNotFound,
			expectedBody: "",
			wantErr:      true,
		},
		{
			name:          "json counter POST - OK",
			request:       "/update/",
			method:        http.MethodPost,
			body:          `{"id":"some_name","type":"counter","delta":5}`,
			expectedCode:  http.StatusOK,
			expectedBody:  `{"id":"some_name","type":"counter","delta":5}`,
			expectedZBody: `{"id":"some_name","type":"counter","delta":15}`,
			wantErr:       false,
		},
		{
			name:          "json second counter POST - OK",
			request:       "/update/",
			method:        http.MethodPost,
			body:          `{"id":"some_name","type":"counter","delta":5}`,
			expectedCode:  http.StatusOK,
			expectedBody:  `{"id":"some_name","type":"counter","delta":10}`,
			expectedZBody: `{"id":"some_name","type":"counter","delta":20}`,
			wantErr:       false,
		},
		{
			name:          "json gauge POST - OK",
			request:       "/update/",
			method:        http.MethodPost,
			body:          `{"id":"some_name","type":"gauge","value":5}`,
			expectedCode:  http.StatusOK,
			expectedBody:  `{"id":"some_name","type":"gauge","value":5}`,
			expectedZBody: `{"id":"some_name","type":"gauge","value":5}`,
			wantErr:       false,
		},
		{
			name:         "json gauge POST - unknown type",
			request:      "/update/",
			method:       http.MethodPost,
			body:         `{"id":"some_name","type":"unknown","value":5}`,
			expectedCode: http.StatusNotImplemented,
			expectedBody: "",
			wantErr:      true,
		},
	}
	tInfo := struct {
		name         string
		request      string
		method       string
		expectedCode int
	}{
		name:         "GET - OK",
		request:      "/",
		method:       http.MethodGet,
		expectedCode: http.StatusOK,
	}

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

	t.Run(tInfo.name, func(t *testing.T) {
		infoResp, _ := testRequest(t, ts, tInfo.method, tInfo.request, "")
		defer infoResp.Body.Close()
		assert.Equal(t, tInfo.expectedCode, infoResp.StatusCode)
	})

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

	t.Run(tInfo.name, func(t *testing.T) {
		infoResp, _ := testGzRequest(t, ts, tInfo.method, tInfo.request, "")
		defer infoResp.Body.Close()
		assert.Equal(t, tInfo.expectedCode, infoResp.StatusCode)
	})

}
