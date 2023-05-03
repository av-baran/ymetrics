package server

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/av-baran/ymetrics/internal/repository/memstor"
	"github.com/av-baran/ymetrics/pkg/interrors"
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

func TestNew(t *testing.T) {
	s := memstor.New()
	assert.NotEmpty(t, New(s, defaultCfg))
}

func TestServer(t *testing.T) {
	repo := memstor.New()
	serv := New(repo, defaultCfg)
	ts := httptest.NewServer(serv.Router)
	defer ts.Close()

	tests := []struct {
		name         string
		request      string
		method       string
		expectedCode int
	}{
		{
			name:         "GET",
			request:      "/update/gauge/name/1.0",
			method:       http.MethodGet,
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "PUT",
			request:      "/update/gauge/name/1.0",
			method:       http.MethodPut,
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "DELETE",
			request:      "/update/gauge/name/1.0",
			method:       http.MethodDelete,
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "POST - Bad Type",
			request:      "/update/unknowntype/name/1.0",
			method:       http.MethodPost,
			expectedCode: http.StatusNotImplemented,
		},
		{
			name:         "POST - Post Bad value",
			request:      "/update/gauge/name/fuuuu",
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "POST - OK",
			request:      "/update/gauge/name/1.01",
			method:       http.MethodPost,
			expectedCode: http.StatusOK,
		},
		{
			name:         "POST - OK",
			request:      "/update/counter/name/1",
			method:       http.MethodPost,
			expectedCode: http.StatusOK,
		},
	}
	for _, v := range tests {
		resp, _ := testRequest(t, ts, v.method, v.request, "")
		defer resp.Body.Close()
		assert.Equalf(t, v.expectedCode, resp.StatusCode, "request: %v. want: %v, got: %v", v.request, v.expectedCode, resp.StatusCode)
	}
}

func TestServerGetValue(t *testing.T) {
	repo := memstor.New()
	serv := New(repo, defaultCfg)
	ts := httptest.NewServer(serv.Router)
	defer ts.Close()

	tests := []struct {
		name         string
		request      string
		method       string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "POST - OK",
			request:      "/update/gauge/name/1.01",
			method:       http.MethodPost,
			expectedCode: http.StatusOK,
			expectedBody: "",
		},
		{
			name:         "POST - OK",
			request:      "/update/counter/name/1",
			method:       http.MethodPost,
			expectedCode: http.StatusOK,
			expectedBody: "",
		},
		{
			name:         "GET",
			request:      "/value/gauge/name",
			method:       http.MethodGet,
			expectedCode: http.StatusOK,
			expectedBody: "1.01",
		},
		{
			name:         "GET",
			request:      "/value/counter/name",
			method:       http.MethodGet,
			expectedCode: http.StatusOK,
			expectedBody: "1",
		},
		{
			name:         "GET",
			request:      "/value/gauge/unknownname",
			method:       http.MethodGet,
			expectedCode: http.StatusNotFound,
			expectedBody: "cannot get gauge metric: " + interrors.ErrMetricNotFound.Error() + "\n",
		},
	}
	for _, v := range tests {
		resp, got := testRequest(t, ts, v.method, v.request, "")
		defer resp.Body.Close()
		assert.Equal(t, v.expectedCode, resp.StatusCode)
		assert.Equal(t, v.expectedBody, got)
	}
}

func TestServerGetAll(t *testing.T) {
	repo := memstor.New()
	serv := New(repo, defaultCfg)
	ts := httptest.NewServer(serv.Router)
	defer ts.Close()

	tests := []struct {
		name         string
		request      string
		method       string
		expectedCode int
	}{
		{
			name:         "GET - OK",
			request:      "/",
			method:       http.MethodGet,
			expectedCode: http.StatusOK,
		},
	}
	for _, v := range tests {
		resp, _ := testRequest(t, ts, v.method, v.request, "")
		defer resp.Body.Close()
		assert.Equal(t, v.expectedCode, resp.StatusCode)
	}
}

func TestServerUpdateJSON(t *testing.T) {
	repo := memstor.New()
	serv := New(repo, defaultCfg)
	ts := httptest.NewServer(serv.Router)
	defer ts.Close()

	tests := []struct {
		name         string
		request      string
		method       string
		body         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "counter POST - OK",
			request:      "/update/",
			method:       http.MethodPost,
			body:         `{"id":"some_name","type":"counter","delta":5}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"id":"some_name","type":"counter","delta":5}`,
		},
		{
			name:         "second counter POST - OK",
			request:      "/update/",
			method:       http.MethodPost,
			body:         `{"id":"some_name","type":"counter","delta":5}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"id":"some_name","type":"counter","delta":10}`,
		},
		{
			name:         "gauge POST - OK",
			request:      "/update/",
			method:       http.MethodPost,
			body:         `{"id":"some_name","type":"gauge","value":5}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"id":"some_name","type":"gauge","value":5}`,
		},
		{
			name:         "gauge POST - unknown type",
			request:      "/update/",
			method:       http.MethodPost,
			body:         `{"id":"some_name","type":"unknown","value":5}`,
			expectedCode: http.StatusNotImplemented,
			expectedBody: "cannot set metric: invalid metric type" + "\n",
		},
	}
	for _, v := range tests {
		resp, got := testRequest(t, ts, v.method, v.request, v.body)
		defer resp.Body.Close()
		assert.Equal(t, v.expectedCode, resp.StatusCode)
		assert.Equal(t, v.expectedBody, got)
	}
}

func TestServerGetJSON(t *testing.T) {
	repo := memstor.New()
	serv := New(repo, defaultCfg)
	ts := httptest.NewServer(serv.Router)
	defer ts.Close()

	tests := []struct {
		name         string
		request      string
		method       string
		body         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "gauge POST - OK",
			request:      "/update/",
			method:       http.MethodPost,
			body:         `{"id":"some_name","type":"gauge","value":5}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"id":"some_name","type":"gauge","value":5}`,
		},
		{
			name:         "gauge GET - OK",
			request:      "/value/",
			method:       http.MethodPost,
			body:         `{"id":"some_name","type":"gauge"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"id":"some_name","type":"gauge","value":5}`,
		},
		{
			name:         "counter POST - OK",
			request:      "/update/",
			method:       http.MethodPost,
			body:         `{"id":"some_name","type":"counter","delta":5}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"id":"some_name","type":"counter","delta":5}`,
		},
		{
			name:         "second counter POST - OK",
			request:      "/update/",
			method:       http.MethodPost,
			body:         `{"id":"some_name","type":"counter","delta":5}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"id":"some_name","type":"counter","delta":10}`,
		},
		{
			name:         "counter GET - OK",
			request:      "/value/",
			method:       http.MethodPost,
			body:         `{"id":"some_name","type":"counter"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"id":"some_name","type":"counter","delta":10}`,
		},
	}
	for _, v := range tests {
		resp, got := testRequest(t, ts, v.method, v.request, v.body)
		defer resp.Body.Close()
		assert.Equal(t, v.expectedCode, resp.StatusCode)
		assert.Equal(t, v.expectedBody, got)
	}
}

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
