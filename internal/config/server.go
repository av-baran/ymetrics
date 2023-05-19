package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	serverDefaultServerAddress   = "localhost:8080"
	serverDefaultLogLevel        = "debug"
	serverDefaultStoreInterval   = "300s"
	serverDefaultStoragePath     = "/tmp/metrics-db.json"
	serverDefaultRestore         = true
	serverDefaultShutdownTimeout = 10 * time.Second

	storageDefaultRequestTimeout = 1 * time.Second
)

type ServerConfig struct {
	ServerAddress string
	StorageConfig StorageConfig
	LoggerConfig  LoggerConfig

	StoreInterval   time.Duration
	FileStoragePath string
	Restore         bool
	ShutdownTimeout time.Duration
}

func NewServerConfig() (*ServerConfig, error) {
	cfg := &ServerConfig{
		StorageConfig: StorageConfig{
			RequestTimeout: storageDefaultRequestTimeout,
		},
		LoggerConfig:    LoggerConfig{},
		ShutdownTimeout: serverDefaultShutdownTimeout,
	}

	parseServerFlags(cfg)

	if a, ok := os.LookupEnv("ADDRESS"); ok {
		cfg.ServerAddress = a
	}
	if l, ok := os.LookupEnv("LOG_LEVEL"); ok {
		cfg.LoggerConfig.Level = l
	}

	if i, ok := os.LookupEnv("STORE_INTERVAL"); ok {
		v, err := str2Interval(i)
		if err != nil {
			return nil, fmt.Errorf("cannot parse config from env (STORE_INTERVAL): %w", err)
		}
		cfg.StoreInterval = v
	}

	if p, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		cfg.FileStoragePath = p
	}

	if l, ok := os.LookupEnv("RESTORE"); ok {
		switch l {
		case "true", "True", "TRUE", "1":
			cfg.Restore = true
		case "false", "False", "FALSE", "0":
			cfg.Restore = false
		default:
			return nil, fmt.Errorf("cannot parse config from env (RESTORE): wrong value")
		}
	}

	if d, ok := os.LookupEnv("DATABASE_DSN"); ok {
		cfg.StorageConfig.DatabaseDSN = d
	}

	return cfg, nil
}

func parseServerFlags(cfg *ServerConfig) error {
	flag.StringVar(&cfg.ServerAddress, "a", serverDefaultServerAddress, "server address and port to listen")
	flag.StringVar(&cfg.LoggerConfig.Level, "l", serverDefaultLogLevel, "log level")
	flag.StringVar(&cfg.FileStoragePath, "f", serverDefaultStoragePath, "file storage path")
	flag.BoolVar(&cfg.Restore, "r", serverDefaultRestore, "restore from file")
	flag.StringVar(&cfg.StorageConfig.DatabaseDSN, "d", "", "database connection string")

	// Переделал потому что в тестах 10-го инкремента в параметрах передают строку "10s" и заранее
	// неизвестно, что будет передано во флаге, строка подходящяя для time.Duration или просто целое.
	var storeIntervalFlag string
	flag.StringVar(&storeIntervalFlag, "i", serverDefaultStoreInterval, "save to file interval")
	s, err := str2Interval(storeIntervalFlag)
	if err != nil {
		return fmt.Errorf("cannot parse store interval flag: %w", err)
	}
	cfg.StoreInterval = s
	return nil
}

func str2Interval(s string) (time.Duration, error) {
	const parseError = "time: missing unit in duration"

	interval, err := time.ParseDuration(s)
	if err == nil {
		return interval, err
	} else if !strings.Contains(err.Error(), parseError) {
		return 0, fmt.Errorf("cannot parse string to interval: %w", err)
	}
	log.Printf(`%s, trying to parse string as integer number of seconds`, err)

	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("cannot parse string to interval: %w", err)
	}

	if i < 0 {
		return 0, fmt.Errorf("interval shouldn't be negative: %w", err)
	}
	interval = time.Duration(i) * time.Second

	return interval, nil
}

func (c *ServerConfig) IsValidStoreFile() error {
	if _, err := os.Stat(c.FileStoragePath); err != nil {
		return fmt.Errorf("backup file is not found: %w", err)
	}
	return nil
}
