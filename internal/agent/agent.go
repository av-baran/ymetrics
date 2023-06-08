package agent

import (
	"context"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/av-baran/ymetrics/internal/config"
	"github.com/av-baran/ymetrics/internal/logger"
	"github.com/av-baran/ymetrics/internal/metric"
	"github.com/go-resty/resty/v2"
)

type Agent struct {
	cfg         *config.AgentConfig
	pollCount   int64
	client      *resty.Client
	pollCounter *pollCounter
}

var randSrc = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
var metricsC = make(chan []metric.Metric, 1)
var errorC = make(chan error, 1)

func NewAgent(cfg *config.AgentConfig) *Agent {
	a := &Agent{cfg, 0, resty.New(), newPollCounter()}
	a.client.SetTimeout(cfg.GetRequestTimeout())
	return a
}

func (a *Agent) Run() {
	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())

	go a.collectMemStats(ctx, &wg)
	go a.collectSysStats(ctx, &wg)
	go a.batchDump(ctx, &wg)

	running := true
	for running {
		select {
		case <-exitSignal:
			running = false
		case err := <-errorC:
			logger.Error("goroutine send error: %s", err)
			running = false
		}
	}
	cancel()
	wg.Wait()
}

func (a *Agent) Shutdown() {
	close(metricsC)
	close(errorC)
}
