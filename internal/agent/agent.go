package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/av-baran/ymetrics/internal/config"
	"github.com/av-baran/ymetrics/internal/metric"
	"github.com/go-resty/resty/v2"
)

type Agent struct {
	cfg       *config.AgentConfig
	pollCount int64
	client    *resty.Client
}

var randSrc = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
var collectedMetrics = []metric.Metric{}

func NewAgent(cfg *config.AgentConfig) *Agent {
	a := &Agent{cfg, 0, resty.New()}
	a.client.SetTimeout(config.RequestTimeout)
	return a
}

func (a *Agent) Run(cfg *config.AgentConfig) {
	pollTicker := time.NewTicker(a.cfg.GetPollInterval())
	defer pollTicker.Stop()
	reportTicker := time.NewTicker(a.cfg.GetReportInterval())
	defer reportTicker.Stop()

	for {
		select {
		case <-pollTicker.C:
			a.collectMetrics()
		case <-reportTicker.C:
			if err := a.dump(); err != nil {
				log.Printf("cannot dump metrics to server: %s", err)
			}
		}
	}
}

func (a *Agent) collectMetrics() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	a.pollCount++

	collectedMetrics = []metric.Metric{
		{
			Name:  "Alloc",
			Value: m.Alloc,
			Type:  metric.GaugeType,
		},
		{
			Name:  "BuckHashSys",
			Value: m.BuckHashSys,
			Type:  metric.GaugeType,
		},
		{
			Name:  "Frees",
			Value: m.Frees,
			Type:  metric.GaugeType,
		},
		{
			Name:  "GCCPUFraction",
			Value: m.GCCPUFraction,
			Type:  metric.GaugeType,
		},
		{
			Name:  "GCSys",
			Value: m.GCSys,
			Type:  metric.GaugeType,
		},
		{
			Name:  "HeapAlloc",
			Value: m.HeapAlloc,
			Type:  metric.GaugeType,
		},
		{
			Name:  "HeapIdle",
			Value: m.HeapIdle,
			Type:  metric.GaugeType,
		},
		{
			Name:  "HeapInuse",
			Value: m.HeapInuse,
			Type:  metric.GaugeType,
		},
		{
			Name:  "HeapObjects",
			Value: m.HeapObjects,
			Type:  metric.GaugeType,
		},
		{
			Name:  "HeapReleased",
			Value: m.HeapReleased,
			Type:  metric.GaugeType,
		},
		{
			Name:  "HeapSys",
			Value: m.HeapSys,
			Type:  metric.GaugeType,
		},
		{
			Name:  "LastGC",
			Value: m.LastGC,
			Type:  metric.GaugeType,
		},
		{
			Name:  "Lookups",
			Value: m.Lookups,
			Type:  metric.GaugeType,
		},
		{
			Name:  "MCacheInuse",
			Value: m.MCacheInuse,
			Type:  metric.GaugeType,
		},
		{
			Name:  "MCacheSys",
			Value: m.MCacheSys,
			Type:  metric.GaugeType,
		},
		{
			Name:  "MSpanInuse",
			Value: m.MSpanInuse,
			Type:  metric.GaugeType,
		},
		{
			Name:  "MSpanSys",
			Value: m.MSpanSys,
			Type:  metric.GaugeType,
		},
		{
			Name:  "Mallocs",
			Value: m.Mallocs,
			Type:  metric.GaugeType,
		},
		{
			Name:  "NextGC",
			Value: m.NextGC,
			Type:  metric.GaugeType,
		},
		{
			Name:  "NumForcedGC",
			Value: m.NumForcedGC,
			Type:  metric.GaugeType,
		},
		{
			Name:  "NumGC",
			Value: m.NumGC,
			Type:  metric.GaugeType,
		},
		{
			Name:  "OtherSys",
			Value: m.OtherSys,
			Type:  metric.GaugeType,
		},
		{
			Name:  "PauseTotalNs",
			Value: m.PauseTotalNs,
			Type:  metric.GaugeType,
		},
		{
			Name:  "StackInuse",
			Value: m.StackInuse,
			Type:  metric.GaugeType,
		},
		{
			Name:  "StackSys",
			Value: m.StackSys,
			Type:  metric.GaugeType,
		},
		{
			Name:  "Sys",
			Value: m.Sys,
			Type:  metric.GaugeType,
		},
		{
			Name:  "TotalAlloc",
			Value: m.TotalAlloc,
			Type:  metric.GaugeType,
		},
		{
			Name:  "PollCount",
			Value: a.pollCount,
			Type:  metric.CounterType,
		},
		{
			Name:  "RandomValue",
			Value: randSrc.Float64(),
			Type:  metric.GaugeType,
		},
	}
}

func (a *Agent) dump() error {
	defer func() { a.pollCount = 0 }()
	for _, m := range collectedMetrics {
		if err := a.sendMetricJSON(m); err != nil {
			return fmt.Errorf("cannot send metric: %w", err)
		}
	}
	return nil
}

func (a *Agent) sendMetricJSON(m metric.Metric) error {
	var buf bytes.Buffer

	encoder := json.NewEncoder(&buf)

	mNew := &metric.Metrics{
		ID:    m.Name,
		MType: string(m.Type),
	}

	switch m.Type {
	case metric.GaugeType:
		v, err := gauge2float64(m.Value)
		if err != nil {
			return fmt.Errorf("error converting metric value to float: %w", err)
		}
		mNew.Value = &v
		log.Printf("metric name: %v, type: %v, val: %v", mNew.ID, mNew.MType, *mNew.Value)
	case metric.CounterType:
		v, ok := m.Value.(int64)
		if !ok {
			return fmt.Errorf("error converting metric value to int")
		}
		mNew.Delta = &v
		log.Printf("metric name: %v, type: %v, delta: %v", mNew.ID, mNew.MType, *mNew.Delta)
	default:
		return errors.New("type not implemented")
	}

	encoder.Encode(&mNew)

	gzBuf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(gzBuf)
	if _, err := zb.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("error while compressing body")
	}
	if err := zb.Close(); err != nil {
		return fmt.Errorf("error while closing gz buffer")
	}

	resp, err := a.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(gzBuf.Bytes()).
		Post(a.cfg.GetURL() + "/update/")

	if err != nil {
		return fmt.Errorf("cannot sent request; resty error: %w", err)
	}

	if resp.StatusCode() >= 300 {
		return fmt.Errorf("remote server respond with no 200 status code: %v", resp.StatusCode())
	}

	return nil
}

func gauge2float64(v interface{}) (float64, error) {
	var res float64

	switch i := v.(type) {
	case uint32:
		res = float64(i)
	case uint64:
		res = float64(i)
	case float64:
		res = float64(i)
	default:
		return 0, fmt.Errorf("cannot convert to float64, type is not implemented")
	}
	return res, nil
}
