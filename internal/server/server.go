package server

import (
	"github.com/av-baran/ymetrics/internal/logger"
	"github.com/av-baran/ymetrics/internal/repository"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	Storage repository.Storager
	Router  *chi.Mux
	Cfg     *ServerConfig
}

func New(s repository.Storager, cfg *ServerConfig) *Server {
	router := chi.NewRouter()
	router.Use(
		gzMiddleware,
		logger.RequestLogMiddlware,
		logger.ResponseLogMiddleware,
	)
	server := &Server{Storage: s, Router: router, Cfg: cfg}
	server.registerRoutes()
	return server
}

func (s *Server) registerRoutes() {
	s.Router.Post("/update/{type}/{name}/{value}", s.UpdateMetricHandler)
	s.Router.Get("/value/{type}/{name}", s.GetMetricHandler)
	s.Router.Get("/", s.GetAllMetricsHandler)

	s.Router.Post("/update/", s.UpdateMetricJSONHandler)
	s.Router.Post("/value/", s.GetMetricJSONHandler)
}
