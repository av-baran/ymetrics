package httphandlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(s Service) chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Route("/", func(r chi.Router) {
		r.Route("/update", func(r chi.Router) {
			r.Post("/{type}/{name}/{value}", UpdateMetricHandler(s))
			r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusMethodNotAllowed)
				w.Write([]byte("Method not allowed"))
			})
		})
		r.Route("/value", func(r chi.Router) {
			r.Get("/{type}/{name}", GetMetricHandler(s))
			r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusMethodNotAllowed)
				w.Write([]byte("Method not allowed"))
			})
		})
		r.Get("/", GetAllMetricsHandler(s))
		r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
		})
	})
	return router
}
