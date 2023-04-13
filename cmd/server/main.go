package main

import (
	"log"
	"net/http"

	"github.com/av-baran/ymetrics/internal/delivery/router"
	"github.com/av-baran/ymetrics/internal/repository/memstorv2"
	"github.com/av-baran/ymetrics/internal/usecase/service"
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
