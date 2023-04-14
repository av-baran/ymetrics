package main

import (
	"log"
	"os"
	"strconv"
	"time"
)

const (
	defaultProtocol = "http://"
)

type Config struct {
	ServerAddress string
	Protocol      string
	URL           string

	PollInterval   time.Duration
	ReportInterval time.Duration
}

type EnvConfig struct {
	ServerAddress  string
	ReportInterval int
	PollInterval   int
}

func getEnvConfig() *EnvConfig {
	e := &EnvConfig{ServerAddress: os.Getenv("ADDRESS")}

	if r, err := strconv.Atoi(os.Getenv("REPORT_INTERVAL")); err != nil {
		log.Printf("Invalid value in ENV variable REPORT_INTERVAL=%v", r)
	} else {
		e.ReportInterval = r
	}

	if p, err := strconv.Atoi(os.Getenv("POLL_INTERVAL")); err != nil {
		log.Printf("Invalid value in ENV variable POLL_INTERVAL=%v", p)
	} else {
		e.PollInterval = p
	}
	return e
}

func NewConfig() *Config {
	var serverAddress string
	var pollInterval int
	var reportInterval int

	flags := parseFlags()
	envs := getEnvConfig()

	if envs.ServerAddress != "" {
		serverAddress = envs.ServerAddress
		log.Printf("Set ADDRESS=%v from ENV", envs.ServerAddress)
	} else {
		serverAddress = flags.ServerAddress
	}

	if envs.ReportInterval != 0 {
		reportInterval = envs.ReportInterval
		log.Printf("Set REPORT_INTERVAL=%v from ENV", envs.ReportInterval)
	} else {
		reportInterval = flags.ReportInterval
	}

	if envs.PollInterval != 0 {
		pollInterval = envs.PollInterval
		log.Printf("Set POLL_INTERVAL=%v from ENV", envs.PollInterval)
	} else {
		pollInterval = flags.PollInterval
	}

	c := &Config{
		ServerAddress:  serverAddress,
		Protocol:       defaultProtocol,
		ReportInterval: time.Duration(reportInterval) * time.Second,
		PollInterval:   time.Duration(pollInterval) * time.Second,
	}
	c.URL = c.Protocol + c.ServerAddress
	return c
}
