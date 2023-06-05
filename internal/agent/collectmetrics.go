package agent

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/av-baran/ymetrics/internal/logger"
	"github.com/av-baran/ymetrics/internal/metric"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/v3/cpu"
)

func (a *Agent) collectMemStats(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	pollTicker := time.NewTicker(a.cfg.GetPollInterval())
	defer pollTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-pollTicker.C:
			collectedMetrics := a.readMemMetrics()
			select {
			case metricsCh <- collectedMetrics:
			default:
				logger.Info("metrics channel is full, cannot send mem stat")
			}
		}
	}
}

func (a *Agent) collectSysStats(ctx context.Context, wg *sync.WaitGroup) {
	pollTicker := time.NewTicker(a.cfg.GetPollInterval())
	defer pollTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-pollTicker.C:
			collectedMetrics, err := a.readSysMetrics()
			if err != nil {
				errorCh <- err
			}
			metricsCh <- collectedMetrics
		}
	}
}

func (a *Agent) readSysMetrics() ([]metric.Metric, error) {
	collectedMetrics := make([]metric.Metric, 0)

	cpuUtil, err := cpu.Percent(0, true)
	if err != nil {
		return nil, fmt.Errorf("cannot get cpu utilization: %w", err)
	}
	for i, v := range cpuUtil {
		cpuMetric := metric.Metric{
			ID:    fmt.Sprintf("CPUutilization%d", i),
			Value: getFloat64Ptr(v),
			MType: metric.GaugeType,
		}
		collectedMetrics = append(collectedMetrics, cpuMetric)
	}

	vm, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("cannot get mem info: %w", err)
	}

	sysMemMetrics := []metric.Metric{
		{
			ID:    "TotalMemory",
			Value: getFloat64Ptr(float64(vm.Total)),
			MType: metric.GaugeType,
		},
		{
			ID:    "FreeMemory",
			Value: getFloat64Ptr(float64(vm.Free)),
			MType: metric.GaugeType,
		},
	}

	collectedMetrics = append(collectedMetrics, sysMemMetrics...)

	return collectedMetrics, nil
}

func (a *Agent) readMemMetrics() []metric.Metric {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	a.pollCounter.add()

	collectedMetrics := []metric.Metric{
		{
			ID:    "Alloc",
			Value: getFloat64Ptr(float64(m.Alloc)),
			MType: metric.GaugeType,
		},
		{
			ID:    "BuckHashSys",
			Value: getFloat64Ptr(float64(m.BuckHashSys)),
			MType: metric.GaugeType,
		},
		{
			ID:    "Frees",
			Value: getFloat64Ptr(float64(m.Frees)),
			MType: metric.GaugeType,
		},
		{
			ID:    "GCCPUFraction",
			Value: getFloat64Ptr(float64(m.GCCPUFraction)),
			MType: metric.GaugeType,
		},
		{
			ID:    "GCSys",
			Value: getFloat64Ptr(float64(m.GCSys)),
			MType: metric.GaugeType,
		},
		{
			ID:    "HeapAlloc",
			Value: getFloat64Ptr(float64(m.HeapAlloc)),
			MType: metric.GaugeType,
		},
		{
			ID:    "HeapIdle",
			Value: getFloat64Ptr(float64(m.HeapIdle)),
			MType: metric.GaugeType,
		},
		{
			ID:    "HeapInuse",
			Value: getFloat64Ptr(float64(m.HeapInuse)),
			MType: metric.GaugeType,
		},
		{
			ID:    "HeapObjects",
			Value: getFloat64Ptr(float64(m.HeapObjects)),
			MType: metric.GaugeType,
		},
		{
			ID:    "HeapReleased",
			Value: getFloat64Ptr(float64(m.HeapReleased)),
			MType: metric.GaugeType,
		},
		{
			ID:    "HeapSys",
			Value: getFloat64Ptr(float64(m.HeapSys)),
			MType: metric.GaugeType,
		},
		{
			ID:    "LastGC",
			Value: getFloat64Ptr(float64(m.LastGC)),
			MType: metric.GaugeType,
		},
		{
			ID:    "Lookups",
			Value: getFloat64Ptr(float64(m.Lookups)),
			MType: metric.GaugeType,
		},
		{
			ID:    "MCacheInuse",
			Value: getFloat64Ptr(float64(m.MCacheInuse)),
			MType: metric.GaugeType,
		},
		{
			ID:    "MCacheSys",
			Value: getFloat64Ptr(float64(m.MCacheSys)),
			MType: metric.GaugeType,
		},
		{
			ID:    "MSpanInuse",
			Value: getFloat64Ptr(float64(m.MSpanInuse)),
			MType: metric.GaugeType,
		},
		{
			ID:    "MSpanSys",
			Value: getFloat64Ptr(float64(m.MSpanSys)),
			MType: metric.GaugeType,
		},
		{
			ID:    "Mallocs",
			Value: getFloat64Ptr(float64(m.Mallocs)),
			MType: metric.GaugeType,
		},
		{
			ID:    "NextGC",
			Value: getFloat64Ptr(float64(m.NextGC)),
			MType: metric.GaugeType,
		},
		{
			ID:    "NumForcedGC",
			Value: getFloat64Ptr(float64(m.NumForcedGC)),
			MType: metric.GaugeType,
		},
		{
			ID:    "NumGC",
			Value: getFloat64Ptr(float64(m.NumGC)),
			MType: metric.GaugeType,
		},
		{
			ID:    "OtherSys",
			Value: getFloat64Ptr(float64(m.OtherSys)),
			MType: metric.GaugeType,
		},
		{
			ID:    "PauseTotalNs",
			Value: getFloat64Ptr(float64(m.PauseTotalNs)),
			MType: metric.GaugeType,
		},
		{
			ID:    "StackInuse",
			Value: getFloat64Ptr(float64(m.StackInuse)),
			MType: metric.GaugeType,
		},
		{
			ID:    "StackSys",
			Value: getFloat64Ptr(float64(m.StackSys)),
			MType: metric.GaugeType,
		},
		{
			ID:    "Sys",
			Value: getFloat64Ptr(float64(m.Sys)),
			MType: metric.GaugeType,
		},
		{
			ID:    "TotalAlloc",
			Value: getFloat64Ptr(float64(m.TotalAlloc)),
			MType: metric.GaugeType,
		},
		{
			ID:    "PollCount",
			Delta: getInt64Ptr(a.pollCounter.get()),
			MType: metric.CounterType,
		},
		{
			ID:    "RandomValue",
			Value: getFloat64Ptr(randSrc.Float64()),
			MType: metric.GaugeType,
		},
	}

	return collectedMetrics
}

func getFloat64Ptr(v float64) *float64 {
	return &v
}

func getInt64Ptr(v int64) *int64 {
	return &v
}
