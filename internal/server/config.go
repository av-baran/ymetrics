package server

import (
	"flag"
	"os"
	"strconv"
	"time"
)

const (
	defaultServerAddress = "localhost:8080"
	defaultLogLevel      = "debug"
	defaultStoreInterval = 300
	defaultStoragePath   = "/tmp/metrics-db.json"
	defaultRestore       = true
)

type ServerConfig struct {
	ServerAddress   string
	LogLevel        string
	StoreInterval   int
	FileStoragePath string
	Restore         bool
}

func NewServerConfig() *ServerConfig {
	cfg := &ServerConfig{}

	parseFlags(cfg)

	if a, ok := os.LookupEnv("ADDRESS"); ok {
		cfg.ServerAddress = a
	}
	if l, ok := os.LookupEnv("LOG_LEVEL"); ok {
		cfg.LogLevel = l
	}

	if i, ok := os.LookupEnv("STORE_INTERVAL"); ok {
		if v, err := strconv.Atoi(i); err == nil {
			cfg.StoreInterval = v
		}
	}

	if p, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		cfg.LogLevel = p
	}

	if l, ok := os.LookupEnv("RESTORE"); ok {
		switch l {
		case "true", "True", "TRUE", "1":
			cfg.Restore = true
		case "false", "False", "FALSE", "0":
			cfg.Restore = false
		}
	}

	return cfg
}

func parseFlags(cfg *ServerConfig) {
	flag.StringVar(&cfg.ServerAddress, "a", defaultServerAddress, "server address and port to listen")
	flag.StringVar(&cfg.LogLevel, "l", defaultLogLevel, "log level")
	flag.IntVar(&cfg.StoreInterval, "i", defaultStoreInterval, "restore interval")
	flag.StringVar(&cfg.FileStoragePath, "f", defaultStoragePath, "file storage path")
	flag.BoolVar(&cfg.Restore, "r", defaultRestore, "restore from file")

	flag.Parse()
}

func (c *ServerConfig) GetStoreInterval() time.Duration {
	return time.Duration(c.StoreInterval) * time.Second
}
