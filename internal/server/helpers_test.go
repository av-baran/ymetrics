package server

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
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
