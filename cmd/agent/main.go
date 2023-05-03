package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/av-baran/ymetrics/internal/agent"
)

func main() {
	cfg := agent.NewAgentConfig()
	a := agent.NewAgent(cfg)
	go a.Run()

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
}
