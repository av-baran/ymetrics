package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/av-baran/ymetrics/internal/repository/memstorv2"
	"github.com/av-baran/ymetrics/internal/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	repo := memstorv2.New()
	serv := server.New(repo)

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Post("/update/{type}/{name}/{value}", serv.UpdateMetricHandler)
	router.Get("/value/{type}/{name}", serv.GetMetricHandler)
	router.Get("/", serv.GetAllMetricsHandler)

	if err := http.ListenAndServe(flagServerAddress, router); err != nil {
		log.Fatalf("cannot run server: %s", err)
	}
}
