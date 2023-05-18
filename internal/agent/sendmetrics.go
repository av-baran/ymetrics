package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/av-baran/ymetrics/internal/metric"
)

func (a *Agent) batchDump() error {
	defer func() { a.pollCount = 0 }()

	if err := a.sendBatchJSON(collectedMetrics); err != nil {
		return fmt.Errorf("cannot send metrics batch: %w", err)
	}
	return nil
}

func (a *Agent) sendBatchJSON(metrics []metric.Metric) error {

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	if err := encoder.Encode(metrics); err != nil {
		return fmt.Errorf("cannot encode metrics: %w", err)
	}

	gzBuf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(gzBuf)
	if _, err := zb.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("error while compressing body: %w", err)
	}
	if err := zb.Close(); err != nil {
		return fmt.Errorf("error while closing gz buffer: %w", err)
	}

	resp, err := a.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(gzBuf).
		// AddRetryCondition().
		Post(a.cfg.GetURL() + "/updates/")
	if err != nil {
		return fmt.Errorf("cannot sent request; resty error: %w", err)
	}

	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusAccepted {
		return fmt.Errorf("remote server respond with unexpected status code: %v", resp.StatusCode())
	}

	return nil
}
