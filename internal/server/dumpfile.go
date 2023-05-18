package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/av-baran/ymetrics/internal/logger"
	"github.com/av-baran/ymetrics/internal/metric"
)

var dumpFileSync = sync.Mutex{}

func (s *Server) restore() error {
	dumpFileSync.Lock()
	defer dumpFileSync.Unlock()

	if err := s.cfg.IsValidStoreFile(); err != nil {
		return fmt.Errorf("cannot restore from backup: %w", err)
	}

	buf, err := os.ReadFile(s.cfg.FileStoragePath)
	if err != nil {
		return fmt.Errorf("cannot read backup file: %w", err)
	}

	metrics := make([]metric.Metric, 0)

	if err := json.Unmarshal(buf, &metrics); err != nil {
		return fmt.Errorf("cannot unmarshal backup file content: %w", err)
	}

	for _, v := range metrics {
		if err := s.Storage.SetMetric(v); err != nil {
			return fmt.Errorf("cannot set metric from backup: %w", err)
		}
	}
	return nil
}

func (s *Server) dumpfile() error {
	dumpFileSync.Lock()
	defer dumpFileSync.Unlock()

	file, err := os.OpenFile(s.cfg.FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("cannot open/create backup file for write: %w", err)
	}
	defer file.Close()

	buf := bufio.NewWriter(file)
	defer buf.Flush()

	encoder := json.NewEncoder(buf)

	metrics, err := s.Storage.GetAllMetrics()
	if err != nil {
		return fmt.Errorf("cannot get all metrics: %w", err)
	}

	if err := encoder.Encode(&metrics); err != nil {
		return fmt.Errorf("cannot encode metrics: %w", err)
	}

	return nil
}

func (s *Server) syncfile() {
	storeInterval := s.cfg.GetStoreInterval()
	if storeInterval <= 0 {
		logger.Info("Store interval is 0. Periodical sync has been disabled.")
		return
	}

	syncTicker := time.NewTicker(storeInterval)
	defer syncTicker.Stop()

	for range syncTicker.C {
		if err := s.dumpfile(); err != nil {
			logger.Fatalf("cannot sync backup file: %s", err)
		}
	}
}
func (s *Server) dumpFileMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		storeInterval := s.cfg.GetStoreInterval()
		if storeInterval == 0 {
			if s.cfg.FileStoragePath == "" {
				return
			}
			s.dumpfile()
		}
	})
}
