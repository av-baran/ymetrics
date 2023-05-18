package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/av-baran/ymetrics/internal/config"
	"github.com/av-baran/ymetrics/internal/logger"
	"github.com/av-baran/ymetrics/internal/repository"
	"github.com/av-baran/ymetrics/internal/repository/memstor"
	"github.com/av-baran/ymetrics/internal/repository/psql"
	"github.com/av-baran/ymetrics/internal/server"
)

func main() {
	cfg, err := config.NewServerConfig()
	if err != nil {
		log.Fatalf("cannot init config: %s", err)
	}

	if err := logger.Init(cfg.LoggerConfig); err != nil {
		log.Fatalf("cannot initialize logger: %s", err)
	}
	defer logger.Sync()

	var repo repository.Storage
	if cfg.StorageConfig.DatabaseDSN != "" {
		repo = psql.New()
		err := repo.Init(cfg.StorageConfig)
		if err != nil {
			logger.Fatalf("cannot init storage: %s", err)
		}
	} else {
		repo = memstor.New()
	}

	srv := server.New(repo, cfg)
	go srv.Run()

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal

	if err := srv.Shutdown(); err != nil {
		logger.Fatalf("cannot gracefully shutdown server: %w", err)
	}

	if err := repo.Shutdown(); err != nil {
		logger.Fatalf("cannot gracefully shutdown storage: %w", err)
	}
}
