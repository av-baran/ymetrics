package config

import (
	"flag"
	"os"
	"strconv"
	"time"
)

const (
	defaultProtocol = "http://"
	RequestTimeout  = time.Second * 1
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
	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "server address and port to listen")
	flag.IntVar(&cfg.ReportInterval, "r", 3, "report interval in seconds")
	flag.IntVar(&cfg.PollInterval, "p", 1, "poll interval in seconds")

	flag.Parse()
}

func (a *AgentConfig) GetURL() string {
	return defaultProtocol + a.ServerAddress
}

func (a *AgentConfig) GetPollInterval() time.Duration {
	return time.Duration(a.PollInterval) * time.Second
}

func (a *AgentConfig) GetReportInterval() time.Duration {
	return time.Duration(a.ReportInterval) * time.Second
}
