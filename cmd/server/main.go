package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/av-baran/ymetrics/internal/repository/memstor"
	"github.com/av-baran/ymetrics/internal/server"
)

var flagServerAddress string

func parseFlags() {
	flag.StringVar(&flagServerAddress, "a", "localhost:8080", "server address and port to listen")
	flag.Parse()

	if a, ok := os.LookupEnv("ADDRESS"); ok {
		flagServerAddress = a
	}

}

func main() {
	parseFlags()

	repo := memstor.New()
	srv := server.New(repo)

	if err := http.ListenAndServe(flagServerAddress, srv.Router); err != nil {
		log.Fatalf("cannot run server: %s", err)
	}
}
