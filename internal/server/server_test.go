package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/av-baran/ymetrics/internal/repository/memstor"
	"github.com/av-baran/ymetrics/pkg/interrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	s := memstor.New()
	assert.NotEmpty(t, New(s))
}

func TestServer(t *testing.T) {
	repo := memstor.New()
	serv := New(repo)
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
		resp, _ := testRequest(t, ts, v.method, v.request)
		defer resp.Body.Close()
		assert.Equal(t, v.expectedCode, resp.StatusCode)
	}
}

func TestServerGetValue(t *testing.T) {
	repo := memstor.New()
	serv := New(repo)
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
			expectedBody: "cannot get gauge metric: " + interrors.ErrMetricNotFound + "\n",
		},
	}
	for _, v := range tests {
		resp, got := testRequest(t, ts, v.method, v.request)
		defer resp.Body.Close()
		assert.Equal(t, v.expectedCode, resp.StatusCode)
		assert.Equal(t, v.expectedBody, got)
	}
}

func TestServerGetAll(t *testing.T) {
	repo := memstor.New()
	serv := New(repo)
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
		resp, _ := testRequest(t, ts, v.method, v.request)
		defer resp.Body.Close()
		assert.Equal(t, v.expectedCode, resp.StatusCode)
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}
