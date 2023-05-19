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
	serverDefaultStoreInterval   = 300
	serverDefaultStoragePath     = "/tmp/metrics-db.json"
	serverDefaultRestore         = true
	serverDefaultShutdownTimeout = time.Second * 10

	storageDefaultRequestTimeout = time.Second * 1
)

type ServerConfig struct {
	ServerAddress string
	StorageConfig StorageConfig
	LoggerConfig  LoggerConfig

	StoreInterval   int
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

	// в тестах 10-го инкремента в параметрах передают строку "2s", до этого целое число
	if i, ok := os.LookupEnv("STORE_INTERVAL"); ok {
		d, err := time.ParseDuration(i)
		if !strings.Contains(err.Error(), "time: missing unit in duration") {
			return nil, fmt.Errorf("cannot parse string to interval: %w", err)
		} else if err != nil {
			log.Printf(`%s, trying to parse string as integer number of seconds`, err)
			v, err := strconv.Atoi(i)
			if err != nil {
				return nil, fmt.Errorf("cannot parse config from env (STORE_INTERVAL): %w", err)
			}
			cfg.StoreInterval = v
		} else {
			cfg.StoreInterval = int(d.Seconds())
		}
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

func parseServerFlags(cfg *ServerConfig) {
	flag.StringVar(&cfg.ServerAddress, "a", serverDefaultServerAddress, "server address and port to listen")
	flag.StringVar(&cfg.LoggerConfig.Level, "l", serverDefaultLogLevel, "log level")
	flag.IntVar(&cfg.StoreInterval, "i", serverDefaultStoreInterval, "save to file interval")
	flag.StringVar(&cfg.FileStoragePath, "f", serverDefaultStoragePath, "file storage path")
	flag.BoolVar(&cfg.Restore, "r", serverDefaultRestore, "restore from file")
	flag.StringVar(&cfg.StorageConfig.DatabaseDSN, "d", "", "database connection string")

	flag.Parse()
}

func (c *ServerConfig) GetStoreInterval() time.Duration {
	return time.Duration(c.StoreInterval) * time.Second
}

func (c *ServerConfig) IsValidStoreFile() error {
	if _, err := os.Stat(c.FileStoragePath); err != nil {
		return fmt.Errorf("backup file is not found: %w", err)
	}
	return nil
}
