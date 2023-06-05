package main

import (
	"log"

	"github.com/av-baran/ymetrics/internal/agent"
	"github.com/av-baran/ymetrics/internal/config"
	"github.com/av-baran/ymetrics/internal/logger"
	"github.com/mailru/easyjson"
)

func main() {

	var x easyjson.Marshaler
	easyjson.Marshal(x)

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

	a.Shutdown()
}
