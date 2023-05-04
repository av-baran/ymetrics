package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/av-baran/ymetrics/internal/metric"
)

var dumpFileSync = sync.Mutex{}

func (s *Server) Restore() error {
	dumpFileSync.Lock()
	defer dumpFileSync.Unlock()

	buf, err := os.ReadFile(s.Cfg.FileStoragePath)
	if err != nil {
		return fmt.Errorf("cannot read backup file: %w", err)
	}

	metrics := make([]metric.Metrics, 0)

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

func (s *Server) Dumpfile() error {
	dumpFileSync.Lock()
	defer dumpFileSync.Unlock()

	file, err := os.OpenFile(s.Cfg.FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("cannot open/create backup file for write: %w", err)
	}
	defer file.Close()

	buf := bufio.NewWriter(file)
	defer buf.Flush()

	encoder := json.NewEncoder(buf)

	metrics := s.Storage.GetAllMetrics()
	if err := encoder.Encode(&metrics); err != nil {
		return fmt.Errorf("cannot encode metrics: %w", err)
	}

	return nil
}

func (s *Server) Syncfile() {
	syncTicker := time.NewTicker(s.Cfg.GetStoreInterval())
	defer syncTicker.Stop()

	for range syncTicker.C {
		if err := s.Dumpfile(); err != nil {
			log.Fatalf("cannot sync backup file: %s", err)
		}
	}
}
