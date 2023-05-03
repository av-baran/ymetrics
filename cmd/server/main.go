package main

import (
	"log"
	"net/http"

	"github.com/av-baran/ymetrics/internal/logger"
	"github.com/av-baran/ymetrics/internal/repository/memstor"
	"github.com/av-baran/ymetrics/internal/server"
)

func main() {
	cfg := server.NewServerConfig()

	repo := memstor.New()
	srv := server.New(repo, cfg)

	if err := logger.Init(srv.Cfg.LogLevel); err != nil {
		log.Fatalf("cannot initialize logger: %s", err)
	}

	if err := http.ListenAndServe(srv.Cfg.ServerAddress, srv.Router); err != nil {
		log.Fatalf("cannot run server: %s", err)
	}
}
