package httphandlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/av-baran/ymetrics/internal/repository/memstorv2"
	"github.com/av-baran/ymetrics/internal/usecase/service"
	"github.com/av-baran/ymetrics/pkg/interrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestNewRouter(t *testing.T) {
	repo := memstorv2.New()
	serv := service.New(repo)
	ts := httptest.NewServer(NewRouter(serv))
	defer ts.Close()

	tests := []struct {
		name         string
		request      string
		method       string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "GET",
			request:      "/update/gauge/name/1.0",
			method:       http.MethodGet,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "Method not allowed",
		},
		{
			name:         "PUT",
			request:      "/update/gauge/name/1.0",
			method:       http.MethodPut,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "Method not allowed",
		},
		{
			name:         "DELETE",
			request:      "/update/gauge/name/1.0",
			method:       http.MethodDelete,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "Method not allowed",
		},
		{
			name:         "POST - Bad Type",
			request:      "/update/unknowntype/name/1.0",
			method:       http.MethodPost,
			expectedCode: http.StatusNotImplemented,
			expectedBody: interrors.ErrInvalidMetricType + "\n",
		},
		{
			name:         "POST - Post Bad value",
			request:      "/update/gauge/name/fuuuu",
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			expectedBody: interrors.ErrInvalidMetricValue + "\n",
		},
		{
			name:         "POST - OK",
			request:      "/update/gauge/name/1.01",
			method:       http.MethodPost,
			expectedCode: http.StatusOK,
			expectedBody: "",
		},
		{
			//FIXME depends on previous test value. Should be same metric name
			name:         "POST - already exists",
			request:      "/update/counter/name/1",
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			expectedBody: interrors.ErrMetricExistsWithAnotherType + "\n",
		},
		{
			//FIXME depends on previous saved value. Should be same metric name
			name:         "GET",
			request:      "/value/gauge/name",
			method:       http.MethodGet,
			expectedCode: http.StatusOK,
			expectedBody: "1.01",
		},
		{
			name:         "GET",
			request:      "/value/gauge/unknownname",
			method:       http.MethodGet,
			expectedCode: http.StatusNotFound,
			expectedBody: interrors.ErrMetricNotFound + "\n",
		},
	}
	for _, v := range tests {
		resp, get := testRequest(t, ts, v.method, v.request)
		defer resp.Body.Close()
		assert.Equal(t, v.expectedCode, resp.StatusCode)
		assert.Equal(t, v.expectedBody, get)
	}
}
