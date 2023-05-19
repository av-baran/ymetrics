package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	agentDefaultProtocol       = "http://"
	agentDefaultServerAddress  = "localhost:8080"
	agentDefaultPollInterval   = 3
	agentDefautlReportInterval = 5
	agentDefaultRequestTimeout = 1
	agentDefaultLogLevel       = "info"
)

type AgentConfig struct {
	ServerAddress string
	LoggerConfig  LoggerConfig

	PollInterval   uint
	ReportInterval uint
	RequestTimeout uint

	SignSecretKey string
}

func NewAgentConfig() (*AgentConfig, error) {
	cfg := &AgentConfig{
		LoggerConfig: LoggerConfig{},
	}

	cfg.RequestTimeout = agentDefaultRequestTimeout

	parseAgentFlags(cfg)

	if a, ok := os.LookupEnv("ADDRESS"); ok {
		cfg.ServerAddress = a
	}

	if v, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		r, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("cannot parse config from env (REPORT_INTERVAL): %w", err)
		}
		cfg.ReportInterval = uint(r)
	}

	if v, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		r, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("cannot parse config from env (POLL_INTERVAL): %w", err)
		}
		cfg.PollInterval = uint(r)
	}

	if k, ok := os.LookupEnv("KEY"); ok {
		cfg.SignSecretKey = k
	}

	return cfg, nil
}

func parseAgentFlags(cfg *AgentConfig) {
	flag.StringVar(&cfg.ServerAddress, "a", agentDefaultServerAddress, "server address and port to listen")
	flag.UintVar(&cfg.ReportInterval, "r", agentDefautlReportInterval, "report interval in seconds")
	flag.UintVar(&cfg.PollInterval, "p", agentDefaultPollInterval, "poll interval in seconds")
	flag.StringVar(&cfg.LoggerConfig.Level, "l", agentDefaultLogLevel, "log level")

	flag.StringVar(&cfg.SignSecretKey, "k", "", "enable data signing")

	flag.Parse()
}

func (a *AgentConfig) GetURL() string {
	return agentDefaultProtocol + a.ServerAddress
}

func (a *AgentConfig) GetPollInterval() time.Duration {
	return time.Duration(a.PollInterval) * time.Second
}

func (a *AgentConfig) GetReportInterval() time.Duration {
	return time.Duration(a.ReportInterval) * time.Second
}

func (a *AgentConfig) GetRequestTimeout() time.Duration {
	return time.Duration(a.RequestTimeout) * time.Second
}
