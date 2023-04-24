package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/av-baran/ymetrics/internal/logger"
	"github.com/av-baran/ymetrics/internal/repository/memstor"
	"github.com/av-baran/ymetrics/internal/server"
)

var (
	flagServerAddress string
	flagLogLevel      string
)

func parseFlags() {
	flag.StringVar(&flagServerAddress, "a", "localhost:8080", "server address and port to listen")
	flag.StringVar(&flagLogLevel, "l", "debug", "log level")
	flag.Parse()

	if a, ok := os.LookupEnv("ADDRESS"); ok {
		flagServerAddress = a
	}
	if l, ok := os.LookupEnv("LOG_LEVEL"); ok {
		flagLogLevel = l
	}
}

func main() {
	parseFlags()

	if err := logger.Init(flagLogLevel); err != nil {
		log.Fatalf("cannot initialize logger: %s", err)
	}

	repo := memstor.New()
	srv := server.New(repo)

	if err := http.ListenAndServe(flagServerAddress, srv.Router); err != nil {
		log.Fatalf("cannot run server: %s", err)
	}
}
