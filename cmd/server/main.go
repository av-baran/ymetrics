package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/av-baran/ymetrics/internal/config"
	"github.com/av-baran/ymetrics/internal/logger"
	"github.com/av-baran/ymetrics/internal/repository"
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

	repo, err := repository.New(cfg.StorageConfig)
	if err != nil {
		logger.Fatalf("cannot create new repository: %s", err)
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
