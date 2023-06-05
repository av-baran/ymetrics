package agent

import (
	"math/rand"
	"os"
	"os/signal"
	"syscall"
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
var doneCh = make(chan struct{})

func NewAgent(cfg *config.AgentConfig) *Agent {
	a := &Agent{cfg, 0, resty.New()}
	a.client.SetTimeout(cfg.GetRequestTimeout())
	return a
}

func (a *Agent) Run() {
	// reportTicker := time.NewTicker(a.cfg.GetReportInterval())
	// defer reportTicker.Stop()
	//
	// running := true
	// for running {
	// 	select {
	// 	case <-pollTicker.C:
	// 		a.collectMetrics()
	// 	case <-reportTicker.C:
	// 		if err := a.batchDump(); err != nil {
	// 			logger.Errorf("cannot dump metrics to server: %s", err)
	// 		}
	// 	case <-done:
	// 		running = false
	// 	}
	// }
	metricsCh := make(chan metric.Metric)
	defer close(metricsCh)

	go a.collectMemStats(doneCh, metricsCh)
	go a.collectSysStats(doneCh, metricsCh)
	go a.batchDump(doneCh, metricsCh)

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
}

func (a *Agent) Shutdown() {
	doneCh <- struct{}{}
}
