package agent

import (
	"math/rand"
	"time"

	"github.com/av-baran/ymetrics/internal/config"
	"github.com/av-baran/ymetrics/internal/logger"
	"github.com/go-resty/resty/v2"
)

type Agent struct {
	cfg       *config.AgentConfig
	pollCount int64
	client    *resty.Client
}

var randSrc = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
var done = make(chan bool)

func NewAgent(cfg *config.AgentConfig) *Agent {
	a := &Agent{cfg, 0, resty.New()}
	a.client.SetTimeout(cfg.GetRequestTimeout())
	return a
}

func (a *Agent) Run() {
	pollTicker := time.NewTicker(a.cfg.GetPollInterval())
	defer pollTicker.Stop()
	reportTicker := time.NewTicker(a.cfg.GetReportInterval())
	defer reportTicker.Stop()

	running := true
	for running {
		select {
		case <-pollTicker.C:
			a.collectMetrics()
		case <-reportTicker.C:
			if err := a.batchDump(); err != nil {
				logger.Errorf("cannot dump metrics to server: %s", err)
			}
		case <-done:
			running = false
		}
	}
}

func (a *Agent) Shutdown() {
	done <- true
}
