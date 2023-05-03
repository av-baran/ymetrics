package agent

import (
	"flag"
	"os"
	"strconv"
	"time"
)

const (
	defaultProtocol       = "http://"
	defaultServerAddress  = "localhost:8080"
	defaultPollInterval   = 3
	defautlReportInterval = 5
	RequestTimeout        = time.Second * 1
)

type AgentConfig struct {
	ServerAddress string

	PollInterval   int
	ReportInterval int
}

func NewAgentConfig() *AgentConfig {
	cfg := &AgentConfig{}

	parseFlags(cfg)

	if a, ok := os.LookupEnv("ADDRESS"); ok {
		cfg.ServerAddress = a
	}

	if v, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		if r, err := strconv.Atoi(v); err == nil {
			cfg.ReportInterval = r
		}
	}

	if v, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		if r, err := strconv.Atoi(v); err == nil {
			cfg.PollInterval = r
		}
	}
	return cfg
}

func parseFlags(cfg *AgentConfig) {
	flag.StringVar(&cfg.ServerAddress, "a", defaultServerAddress, "server address and port to listen")
	flag.IntVar(&cfg.ReportInterval, "r", defautlReportInterval, "report interval in seconds")
	flag.IntVar(&cfg.PollInterval, "p", defaultPollInterval, "poll interval in seconds")

	flag.Parse()
}

func (a *AgentConfig) getURL() string {
	return defaultProtocol + a.ServerAddress
}

func (a *AgentConfig) getPollInterval() time.Duration {
	return time.Duration(a.PollInterval) * time.Second
}

func (a *AgentConfig) getReportInterval() time.Duration {
	return time.Duration(a.ReportInterval) * time.Second
}
