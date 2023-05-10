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
	Storage    repository.Storager
	Router     *chi.Mux
	cfg        *config.ServerConfig
	httpServer *http.Server
}

func New(s repository.Storager, cfg *config.ServerConfig) *Server {
	router := chi.NewRouter()
	router.Use(
		gzMiddleware,
		logger.RequestLogMiddlware,
		logger.ResponseLogMiddleware,
	)
	server := &Server{Storage: s, Router: router, cfg: cfg}
	server.registerRoutes()
	server.httpServer = &http.Server{
		Addr:    server.cfg.ServerAddress,
		Handler: server.Router,
	}

	return server
}

func (s *Server) registerRoutes() {
	s.Router.Get("/", s.GetAllMetricsHandler)

	s.Router.Post("/value/", s.GetMetricJSONHandler)
	s.Router.Get("/value/{type}/{name}", s.GetMetricHandler)

	s.Router.Route("/update", func(r chi.Router) {
		r.Use(s.dumpFileMiddleware)
		r.Post("/", s.UpdateMetricJSONHandler)
		r.Post("/{type}/{name}/{value}", s.UpdateMetricHandler)
	})
}

func (s *Server) Run() {
	if s.cfg.Restore {
		if err := s.restore(); err != nil {
			logger.Fatalf("cannot run server: %s", err)
		}
	}

	if err := s.httpServer.ListenAndServe(); err != nil {
		logger.Fatalf("cannot run server: %s", err)
	}

	// если файла нет, мы создадим его в dumpfile(), который будет вызван в syncfile()
	go s.syncfile()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
	defer cancel()
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("cannot gracefully shutdown server: %w", err)
	}
	s.dumpfile()
	return nil
}
