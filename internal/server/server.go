package server

import (
	"github.com/av-baran/ymetrics/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	Storage repository.Storager
	Router  *chi.Mux
}

func New(s repository.Storager) *Server {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	server := &Server{Storage: s, Router: router}
	server.registerRoutes()
	return server
}

func (s *Server) registerRoutes() {
	s.Router.Post("/update/{type}/{name}/{value}", s.UpdateMetricHandler)
	s.Router.Get("/value/{type}/{name}", s.GetMetricHandler)
	s.Router.Get("/", s.GetAllMetricsHandler)
}