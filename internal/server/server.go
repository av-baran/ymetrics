package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/av-baran/ymetrics/internal/config"
	"github.com/av-baran/ymetrics/internal/logger"
	"github.com/av-baran/ymetrics/internal/repository"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	Storage    repository.Storage
	Router     *chi.Mux
	cfg        *config.ServerConfig
	httpServer *http.Server
}

func New(s repository.Storage, cfg *config.ServerConfig) *Server {
	router := chi.NewRouter()
	server := &Server{Storage: s, Router: router, cfg: cfg}
	server.registerRoutes()
	server.httpServer = &http.Server{
		Addr:    server.cfg.ServerAddress,
		Handler: server.Router,
	}

	return server
}

func (s *Server) registerRoutes() {
	s.Router.Route("/", func(r chi.Router) {
		r.Use(
			s.checkSignMiddleware,
			gzMiddleware,
			logger.RequestLogMiddlware,
			logger.ResponseLogMiddleware,
			s.addSignMiddleware,
		)

		r.Get("/", s.GetAllMetricsHandler)
		r.Get("/ping", s.PingDBHandler)

		r.Route("/update", func(r chi.Router) {
			r.Use(s.dumpFileMiddleware)
			r.Post("/", s.UpdateMetricJSONHandler)
			r.Post("/{type}/{name}/{value}", s.UpdateMetricHandler)
		})

		r.Route("/updates", func(r chi.Router) {
			r.Use(s.dumpFileMiddleware)
			r.Post("/", s.UpdateBatchJSONHandler)
		})

		r.Route("/value", func(r chi.Router) {
			r.Post("/", s.GetMetricJSONHandler)
			r.Get("/{type}/{name}", s.GetMetricHandler)
		})
	})
}

func (s *Server) Run() {
	if s.cfg.Restore {
		if err := s.restore(); err != nil {
			logger.Errorf("cannot restore from backup: %s", err)
		}
	}

	if s.cfg.FileStoragePath != "" {
		go s.syncfile()
	}

	if err := s.httpServer.ListenAndServe(); err != nil {
		logger.Fatalf("cannot run server: %s", err)
	}
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
	defer cancel()
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("cannot gracefully shutdown server: %w", err)
	}

	if s.cfg.FileStoragePath != "" {
		s.dumpfile()
	}

	return nil
}
