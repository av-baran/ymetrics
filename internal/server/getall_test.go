package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/av-baran/ymetrics/internal/repository/memstor"
	"github.com/stretchr/testify/assert"
)

func TestGetAll(t *testing.T) {
	repo := memstor.New()
	serv := New(repo, defaultCfg)
	ts := httptest.NewServer(serv.Router)
	defer ts.Close()

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

	t.Run(tInfo.name, func(t *testing.T) {
		infoResp, _ := testRequest(t, ts, tInfo.method, tInfo.request, "")
		defer infoResp.Body.Close()
		assert.Equal(t, tInfo.expectedCode, infoResp.StatusCode)
	})
}

func TestGetAllWithGzCompression(t *testing.T) {
	repo := memstor.New()
	serv := New(repo, defaultCfg)
	ts := httptest.NewServer(serv.Router)
	defer ts.Close()

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

	t.Run(tInfo.name, func(t *testing.T) {
		infoResp, _ := testGzRequest(t, ts, tInfo.method, tInfo.request, "")
		defer infoResp.Body.Close()
		assert.Equal(t, tInfo.expectedCode, infoResp.StatusCode)
	})
}
