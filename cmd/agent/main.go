package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/av-baran/ymetrics/internal/agent"
	"github.com/av-baran/ymetrics/internal/config"
	"github.com/av-baran/ymetrics/internal/logger"
)

func main() {
	cfg, err := config.NewAgentConfig()
	if err != nil {
		log.Fatalf("cannot init config: %s", err)
	}

	if err := logger.Init(cfg.LoggerConfig); err != nil {
		log.Fatalf("cannot init logger: %s", err)
	}
	defer logger.Sync()

	a := agent.NewAgent(cfg)
	a.Run()

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal

	a.Shutdown()
}
