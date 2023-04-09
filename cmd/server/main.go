package main

import (
	"net/http"

	"github.com/av-baran/ymetrics/internal/handlers"
	storage "github.com/av-baran/ymetrics/internal/storage/mem"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	repo := storage.New()

	mux := http.NewServeMux()
	mux.HandleFunc("/update/", handlers.UpdateMetricHandler(repo))
	return http.ListenAndServe("0.0.0.0:8080", mux)
}
