package router

import (
	"net/http"

	"github.com/av-baran/ymetrics/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(repo handlers.Storage) chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Route("/", func(r chi.Router) {
		r.Route("/update", func(r chi.Router) {
			r.Post("/{type}/{name}/{value}", handlers.UpdateMetricHandler(repo))
			r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusMethodNotAllowed)
				w.Write([]byte("Method not allowed"))
			})
		})
		r.Route("/value", func(r chi.Router) {
			r.Get("/{type}/{name}", handlers.GetMetricHandler(repo))
			r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusMethodNotAllowed)
				w.Write([]byte("Method not allowed"))
			})
		})
	})
	return router
}
