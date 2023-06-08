package agent

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/av-baran/ymetrics/internal/logger"
	"github.com/av-baran/ymetrics/internal/metric"
	"github.com/av-baran/ymetrics/pkg/interrors"
	"github.com/go-resty/resty/v2"
	"golang.org/x/net/context"
)

func (a *Agent) batchDump(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(a.cfg.RateLimit)

	for i := 0; i < a.cfg.RateLimit; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case m := <-metricsC:
					if err := a.sendBatchJSON(m); err != nil {
						logger.Errorf("cannot send metrics batch: %s", err)
					}
					a.pollCounter.reset()
				case <-ctx.Done():
					return
				}
			}
		}()
	}
}

func (a *Agent) sendBatchJSON(metrics []metric.Metric) error {
	headers := map[string]string{
		"Content-Type":     "application/json",
		"Content-Encoding": "gzip",
	}

	buf, err := createRequestJSON(metrics)
	if err != nil {
		return fmt.Errorf("cannot create request body: %w", err)
	}

	if a.cfg.SignSecretKey != "" {
		sign := signBody(a.cfg.SignSecretKey, buf.Bytes())
		headerValue := hex.EncodeToString(sign)
		headers["HashSHA256"] = headerValue
	}

	body, err := compressBuffer(buf)
	if err != nil {
		return fmt.Errorf("cannot compress request body: %w", err)
	}

	var resp *resty.Response
	err = interrors.RetryOnErr(a.cfg.RetryConfig, func() error {
		var restyErr error
		resp, restyErr = a.client.R().
			SetHeaders(headers).
			SetBody(body).
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

func compressBuffer(buf *bytes.Buffer) (*bytes.Buffer, error) {
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

func createRequestJSON(metrics []metric.Metric) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	if err := encoder.Encode(metrics); err != nil {
		return nil, fmt.Errorf("cannot encode metrics: %w", err)
	}
	return &buf, nil
}

func signBody(key string, body []byte) []byte {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(body)
	result := h.Sum(nil)

	return result
}
