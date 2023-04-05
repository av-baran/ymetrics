package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

type MetricKind string

const (
	Gauge   = MetricKind("gauge")
	Counter = MetricKind("counter")

	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
	serverAddress  = "http://localhost:8080"
)

type metric struct {
	name  string
	value interface{}
	kind  MetricKind
}

func main() {
	var m runtime.MemStats
	var pollCount int64

	rand.Seed(time.Now().UTC().UnixNano())

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for range ticker.C {
		runtime.ReadMemStats(&m)
		pollCount++

		metrics := []metric{
			{"Alloc", m.Alloc, Gauge},
			{"BuckHashSys", m.BuckHashSys, Gauge},
			{"Frees", m.Frees, Gauge},
			{"GCCPUFraction", m.GCCPUFraction, Gauge},
			{"GCSys", m.GCSys, Gauge},
			{"HeapAlloc", m.HeapAlloc, Gauge},
			{"HeapIdle", m.HeapIdle, Gauge},
			{"HeapInuse", m.HeapInuse, Gauge},
			{"HeapObjects", m.HeapObjects, Gauge},
			{"HeapReleased", m.HeapReleased, Gauge},
			{"HeapSys", m.HeapSys, Gauge},
			{"LastGC", m.LastGC, Gauge},
			{"Lookups", m.Lookups, Gauge},
			{"MCacheInuse", m.MCacheInuse, Gauge},
			{"MCacheSys", m.MCacheSys, Gauge},
			{"MSpanInuse", m.MSpanInuse, Gauge},
			{"MSpanSys", m.MSpanSys, Gauge},
			{"Mallocs", m.Mallocs, Gauge},
			{"NextGC", m.NextGC, Gauge},
			{"NumForcedGC", m.NumForcedGC, Gauge},
			{"NumGC", m.NumGC, Gauge},
			{"OtherSys", m.OtherSys, Gauge},
			{"PauseTotalNs", m.PauseTotalNs, Gauge},
			{"StackInuse", m.StackInuse, Gauge},
			{"StackSys", m.StackSys, Gauge},
			{"Sys", m.Sys, Gauge},
			{"TotalAlloc", m.TotalAlloc, Gauge},
			{"PollCount", pollCount, Counter},
			{"RandomValue", rand.Float64(), Gauge},
		}

		for _, metric := range metrics {
			if err := sendMetric(metric); err != nil {
				log.Printf(err.Error())
			}
		}
	}
}

func sendMetric(m metric) error {
	switch m.kind {
	case Gauge:
		m.value = fmt.Sprintf("%v", m.value)
	case Counter:
		m.value = fmt.Sprintf("%v", m.value)
	default:
		return errors.New("type not implemented")
	}

	url := fmt.Sprintf("%s/update/%s/%s/%v", serverAddress, m.kind, m.name, m.value)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	defer resp.Body.Close()
	return nil
}
