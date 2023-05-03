package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/av-baran/ymetrics/internal/logger"
	"github.com/av-baran/ymetrics/internal/repository/memstor"
	"github.com/av-baran/ymetrics/internal/server"
	"golang.org/x/net/context"
)

func main() {
	repo := memstor.New()
	cfg := server.NewServerConfig()
	srv := server.New(repo, cfg)

	if err := logger.Init(srv.Cfg.LogLevel); err != nil {
		log.Fatalf("cannot initialize logger: %s", err)
	}

	if srv.Cfg.Restore {
		if _, err := os.Stat(cfg.FileStoragePath); err == nil {
			if err := srv.Restore(); err != nil {
				logger.Log.Sugar().Debugln("error while restoring from file: %s", err.Error())
			}
		} else {
			logger.Log.Sugar().Debugln("backoup file is not exist: %s", err.Error())
		}
	}

	httpServer := &http.Server{Addr: srv.Cfg.ServerAddress, Handler: srv.Router}
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatalf("cannot run server: %s", err)
		}
	}()

	go srv.Syncfile()

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal

	httpServer.Shutdown(context.Background())
	srv.Dumpfile()
}
