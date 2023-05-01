package agent

import (
	"bytes"
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
	return &Agent{cfg, 0, resty.New()}
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

func (a *Agent) sendMetric(m metric.Metric) error {
	switch m.Type {
	case metric.GaugeType, metric.CounterType:
		m.Value = fmt.Sprintf("%v", m.Value)
	default:
		return errors.New("type not implemented")
	}

	resp, err := a.client.R().
		SetHeader("Content-Type", "text/plain").
		SetPathParams(map[string]string{
			"name":  m.Name,
			"type":  string(m.Type),
			"value": m.Value.(string),
		}).
		Post(a.cfg.GetURL() + "/update/{type}/{name}/{value}")

	if err != nil {
		return fmt.Errorf("cannot send metric: %w", err)
	}

	if resp.StatusCode() >= 300 {
		return fmt.Errorf("remote server respond with no 200 status code: %v", resp.StatusCode())
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
		v, _ := m.Value.(float64)
		mNew.Value = &v
	case metric.CounterType:
		v, _ := m.Value.(int64)
		mNew.Delta = &v
	default:
		return errors.New("type not implemented")
	}

	encoder.Encode(&mNew)
	log.Printf("%+v", mNew)

	resp, err := a.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(buf.Bytes()).
		Post(a.cfg.GetURL() + "/update/")

	if err != nil {
		return fmt.Errorf("resty error with statuscode %v: %w", resp.StatusCode(), err)
	}

	if resp.StatusCode() >= 300 {
		return fmt.Errorf("remote server respond with no 200 status code: %v", resp.StatusCode())
	}

	return nil
}
