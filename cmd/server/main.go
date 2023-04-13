package main

import (
	"log"
	"net/http"

	"github.com/av-baran/ymetrics/internal/router"
	"github.com/av-baran/ymetrics/internal/service"
	"github.com/av-baran/ymetrics/internal/storage/memstorv2"
)

func main() {
	parseFlags()

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	repo := memstorv2.New()
	serv := service.New(repo)
	log.Printf("Starting server on %v", flagServerAddress)
	return http.ListenAndServe(flagServerAddress, router.New(serv))
}
