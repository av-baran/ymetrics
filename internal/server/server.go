package server

import (
	"bytes"
	"context"
	"crypto/hmac"
	"fmt"
	"io/ioutil"
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
	s.Router.Get("/ping", s.PingDBHandler)

	s.Router.Post("/value/", s.GetMetricJSONHandler)
	s.Router.Get("/value/{type}/{name}", s.GetMetricHandler)

	s.Router.Route("/update", func(r chi.Router) {
		r.Use(s.dumpFileMiddleware)
		r.Post("/", s.UpdateMetricJSONHandler)
		r.Post("/{type}/{name}/{value}", s.UpdateMetricHandler)
	})
	s.Router.Route("/updates", func(r chi.Router) {
		r.Use(s.dumpFileMiddleware)
		r.Post("/", s.UpdateBatchJSONHandler)
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

func (s *Server) checkSignMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := s.cfg.SignSecretKey
		if key == "" {
			h.ServeHTTP(w, r)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			sendError(w, "cannot read request body", err)
			return
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		if hmac.Equal([]byte(key), body) {
			h.ServeHTTP(w, r)
		} else {
			sendError(w, "cannot close request body", err)
			return
		}
	})
}
