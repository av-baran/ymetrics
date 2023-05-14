package metric

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
)

const (
	GaugeType   = "gauge"
	CounterType = "counter"
	UnknownType = "unknown"
)

type Metric struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func (m *Metric) ToJSON() ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.Encode(&m)

	gzBuf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(gzBuf)
	if _, err := zb.Write(buf.Bytes()); err != nil {
		return nil, fmt.Errorf("error while compressing body: %w", err)
	}
	if err := zb.Close(); err != nil {
		return nil, fmt.Errorf("error while closing gz buffer: %w", err)
	}

	return gzBuf.Bytes(), nil
}
