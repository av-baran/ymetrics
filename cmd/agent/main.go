package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/av-baran/ymetrics/internal/entity/metric"
	"github.com/go-resty/resty/v2"
)

type inMetric struct {
	Name  string
	Value interface{}
	Type  metric.Type
}

var randSrc = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

func main() {
	cfg := NewConfig()
	log.Printf("Starting agent with poll interval: %v, report interval: %v, target server: %v",
		cfg.PollInterval,
		cfg.ReportInterval,
		cfg.ServerAddress,
	)
	run(cfg)
}

func run(c *Config) {
	var pollCount uint64

	pollTicker := time.NewTicker(c.PollInterval)
	defer pollTicker.Stop()
	reportTicker := time.NewTicker(c.ReportInterval)
	defer reportTicker.Stop()

	go func() {
		inputMetrics := collectMetrics(pollCount)
		for {
			select {
			case <-pollTicker.C:
				pollCount++
				inputMetrics = collectMetrics(pollCount)
			case <-reportTicker.C:
				for _, metric := range inputMetrics {
					if err := sendMetric(c.URL, metric); err != nil {
						log.Print(err.Error())
					}
				}
				pollCount = 0
			}
		}
	}()

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
}

func collectMetrics(pollCount uint64) []inMetric {
	var m runtime.MemStats

	runtime.ReadMemStats(&m)

	return []inMetric{
		{"Alloc", m.Alloc, metric.GaugeType},
		{"BuckHashSys", m.BuckHashSys, metric.GaugeType},
		{"Frees", m.Frees, metric.GaugeType},
		{"GCCPUFraction", m.GCCPUFraction, metric.GaugeType},
		{"GCSys", m.GCSys, metric.GaugeType},
		{"HeapAlloc", m.HeapAlloc, metric.GaugeType},
		{"HeapIdle", m.HeapIdle, metric.GaugeType},
		{"HeapInuse", m.HeapInuse, metric.GaugeType},
		{"HeapObjects", m.HeapObjects, metric.GaugeType},
		{"HeapReleased", m.HeapReleased, metric.GaugeType},
		{"HeapSys", m.HeapSys, metric.GaugeType},
		{"LastGC", m.LastGC, metric.GaugeType},
		{"Lookups", m.Lookups, metric.GaugeType},
		{"MCacheInuse", m.MCacheInuse, metric.GaugeType},
		{"MCacheSys", m.MCacheSys, metric.GaugeType},
		{"MSpanInuse", m.MSpanInuse, metric.GaugeType},
		{"MSpanSys", m.MSpanSys, metric.GaugeType},
		{"Mallocs", m.Mallocs, metric.GaugeType},
		{"NextGC", m.NextGC, metric.GaugeType},
		{"NumForcedGC", m.NumForcedGC, metric.GaugeType},
		{"NumGC", m.NumGC, metric.GaugeType},
		{"OtherSys", m.OtherSys, metric.GaugeType},
		{"PauseTotalNs", m.PauseTotalNs, metric.GaugeType},
		{"StackInuse", m.StackInuse, metric.GaugeType},
		{"StackSys", m.StackSys, metric.GaugeType},
		{"Sys", m.Sys, metric.GaugeType},
		{"TotalAlloc", m.TotalAlloc, metric.GaugeType},
		{"PollCount", pollCount, metric.CounterType},
		{"RandomValue", randSrc.Float64(), metric.GaugeType},
	}
}

func sendMetric(srv string, m inMetric) error {
	switch m.Type {
	case metric.GaugeType, metric.CounterType:
		m.Value = fmt.Sprintf("%v", m.Value)
	default:
		return errors.New("type not implemented")
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "text/plain").
		SetPathParams(map[string]string{
			"name":  m.Name,
			"type":  string(m.Type),
			"value": m.Value.(string),
		}).
		Post(srv + "/update/{type}/{name}/{value}")

	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("http status code: %v", resp.StatusCode())
	}
	return nil
}
