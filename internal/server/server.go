package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/av-baran/ymetrics/internal/config"
	"github.com/av-baran/ymetrics/internal/logger"
	"github.com/av-baran/ymetrics/internal/repository"
	"github.com/av-baran/ymetrics/pkg/interrors"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	Storage    repository.Storager
	Router     *chi.Mux
	cfg        *config.ServerConfig
	httpServer *http.Server
	db         *sql.DB
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
	s.Router.Get("/ping/", s.PingDBHandler)

	s.Router.Post("/value/", s.GetMetricJSONHandler)
	s.Router.Get("/value/{type}/{name}", s.GetMetricHandler)

	s.Router.Route("/update", func(r chi.Router) {
		r.Use(s.dumpFileMiddleware)
		r.Post("/", s.UpdateMetricJSONHandler)
		r.Post("/{type}/{name}/{value}", s.UpdateMetricHandler)
	})
}

func (s *Server) Run() {
	if s.cfg.DatabaseDSN != "" {
		err := s.initDB()
		if err != nil {
			logger.Fatalf("cannot run server: %s", err)
		}
	}

	if s.cfg.Restore {
		if err := s.restore(); err != nil {
			logger.Errorf("cannot restore from backup: %s", err)
		}
	}
	go s.syncfile()

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

	s.dumpfile()

	if err := s.db.Close(); err != nil {
		return fmt.Errorf("cannot close DB connection: %w", err)
	}

	return nil
}

func (s *Server) initDB() error {
	db, err := sql.Open("pgx", s.cfg.DatabaseDSN)
	if err != nil {
		return fmt.Errorf("cannot create new DB connection: %w", err)
	}
	s.db = db

	if err := s.pingDB(); err != nil {
		return fmt.Errorf("cannot init DB: %w", err)
	}

	return nil
}

func (s *Server) pingDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("%w: %w", interrors.ErrPingDB, err)
	}

	return nil
}
