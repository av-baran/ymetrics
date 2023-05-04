package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/av-baran/ymetrics/internal/agent"
	"github.com/av-baran/ymetrics/internal/config"
)

func main() {
	cfg := config.NewAgentConfig()

	a := agent.NewAgent(cfg)
	go a.Run(cfg)

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
}
