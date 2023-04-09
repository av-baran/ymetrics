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

	"github.com/av-baran/ymetrics/internal/entities/metric"
	"github.com/go-resty/resty/v2"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
	serverAddress  = "http://localhost:8080"
)

type inMetric struct {
	Name  string
	Value interface{}
	Type  metric.Type
}

func main() {
	var pollCount uint64
	rand.Seed(time.Now().UTC().UnixNano())

	pollTicker := time.NewTicker(pollInterval)
	defer pollTicker.Stop()
	reportTicker := time.NewTicker(reportInterval)
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
					if err := sendMetric(serverAddress, metric); err != nil {
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
		{"Alloc", m.Alloc, metric.Gauge},
		{"BuckHashSys", m.BuckHashSys, metric.Gauge},
		{"Frees", m.Frees, metric.Gauge},
		{"GCCPUFraction", m.GCCPUFraction, metric.Gauge},
		{"GCSys", m.GCSys, metric.Gauge},
		{"HeapAlloc", m.HeapAlloc, metric.Gauge},
		{"HeapIdle", m.HeapIdle, metric.Gauge},
		{"HeapInuse", m.HeapInuse, metric.Gauge},
		{"HeapObjects", m.HeapObjects, metric.Gauge},
		{"HeapReleased", m.HeapReleased, metric.Gauge},
		{"HeapSys", m.HeapSys, metric.Gauge},
		{"LastGC", m.LastGC, metric.Gauge},
		{"Lookups", m.Lookups, metric.Gauge},
		{"MCacheInuse", m.MCacheInuse, metric.Gauge},
		{"MCacheSys", m.MCacheSys, metric.Gauge},
		{"MSpanInuse", m.MSpanInuse, metric.Gauge},
		{"MSpanSys", m.MSpanSys, metric.Gauge},
		{"Mallocs", m.Mallocs, metric.Gauge},
		{"NextGC", m.NextGC, metric.Gauge},
		{"NumForcedGC", m.NumForcedGC, metric.Gauge},
		{"NumGC", m.NumGC, metric.Gauge},
		{"OtherSys", m.OtherSys, metric.Gauge},
		{"PauseTotalNs", m.PauseTotalNs, metric.Gauge},
		{"StackInuse", m.StackInuse, metric.Gauge},
		{"StackSys", m.StackSys, metric.Gauge},
		{"Sys", m.Sys, metric.Gauge},
		{"TotalAlloc", m.TotalAlloc, metric.Gauge},
		{"PollCount", pollCount, metric.Counter},
		{"RandomValue", rand.Float64(), metric.Gauge},
	}
}

func sendMetric(srv string, m inMetric) error {
	switch m.Type {
	case metric.Gauge, metric.Counter:
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
