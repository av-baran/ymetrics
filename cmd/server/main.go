package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/av-baran/ymetrics/internal/config"
	"github.com/av-baran/ymetrics/internal/logger"
	"github.com/av-baran/ymetrics/internal/repository/memstor"
	"github.com/av-baran/ymetrics/internal/server"
)

func main() {
	cfg, err := config.NewServerConfig()
	if err != nil {
		logger.Fatalf("cannot init config: %s", err)
	}

	if err := logger.Init(cfg.LoggerConfig); err != nil {
		log.Fatalf("cannot initialize logger: %s", err)
	}
	defer logger.Sync()

	repo := memstor.New()
	// оставил server.New() чтобы было удобнее тестировать.
	// Если сделать server.Start, в котором вызывать Run() сразу после создания,
	// тогда в тестах не запустить сервер через http.testServer,
	// как это обойти без лишнего усложнения не смог придумать.
	srv := server.New(repo, cfg)
	go srv.Run()

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal

	if err := srv.Shutdown(); err != nil {
		logger.Fatalf("cannot gracefully shutdown server: %w", err)
	}
}
