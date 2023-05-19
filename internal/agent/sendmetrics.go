package agent

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/av-baran/ymetrics/internal/logger"
	"github.com/av-baran/ymetrics/internal/metric"
	"github.com/av-baran/ymetrics/pkg/interrors"
	"github.com/go-resty/resty/v2"
)

func (a *Agent) batchDump() error {
	defer func() { a.pollCount = 0 }()

	if err := a.sendBatchJSON(collectedMetrics); err != nil {
		return fmt.Errorf("cannot send metrics batch: %w", err)
	}
	return nil
}

func (a *Agent) sendBatchJSON(metrics []metric.Metric) error {

	buf, err := createRequestBody(metrics)
	if err != nil {
		return fmt.Errorf("cannot create request body: %w", err)
	}

	headers := map[string]string{
		"Content-Type":     "application/json",
		"Content-Encoding": "gzip",
	}

	if a.cfg.SignSecretKey != "" {
		sign := signBody(a.cfg.SignSecretKey, buf.Bytes())
		headers["HashSHA256"] = sign
	}

	var resp *resty.Response
	err = interrors.RetryOnErr(func() error {
		var restyErr error
		resp, restyErr = a.client.R().
			SetHeaders(headers).
			SetBody(buf).
			Post(a.cfg.GetURL() + "/updates/")
		return restyErr
	})
	if err != nil {
		return fmt.Errorf("cannot sent request; resty error: %w", err)
	}

	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusAccepted {
		return fmt.Errorf("remote server respond with unexpected status code: %v", resp.StatusCode())
	}

	return nil
}

func createRequestBody(metrics []metric.Metric) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	if err := encoder.Encode(metrics); err != nil {
		return nil, fmt.Errorf("cannot encode metrics: %w", err)
	}

	gzBuf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(gzBuf)
	if _, err := zb.Write(buf.Bytes()); err != nil {
		return nil, fmt.Errorf("error while compressing body: %w", err)
	}
	if err := zb.Close(); err != nil {
		return nil, fmt.Errorf("error while closing gz buffer: %w", err)
	}

	return gzBuf, nil
}

func signBody(key string, body []byte) string {

	h := hmac.New(sha256.New, []byte(key))
	h.Write(body)
	result := h.Sum(nil)
	logger.Info("result: %v", string(result))

	return string(result)
}
