package main

import (
	"net/http"

	"github.com/av-baran/ymetrics/internal/handlers"
	storage "github.com/av-baran/ymetrics/internal/storage/mem"
	"github.com/go-chi/chi/v5"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	repo := storage.New()

	router := chi.NewRouter()
	router.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", handlers.UpdateMetricHandler(repo))
		r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
		})
	})
	return http.ListenAndServe("0.0.0.0:8080", router)
}
