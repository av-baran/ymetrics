package main

import (
	"log"
	"net/http"

	"github.com/av-baran/ymetrics/internal/router"
	"github.com/av-baran/ymetrics/internal/storage/memstor"
)

func main() {
	parseFlags()

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	repo := memstor.New()
	log.Printf("Starting server on %v", flagServerAddress)
	return http.ListenAndServe(flagServerAddress, router.New(repo))
}
