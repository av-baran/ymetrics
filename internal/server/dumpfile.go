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

	file, err := os.OpenFile(s.Cfg.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("cannot open/create dump file for write: %w", err)
	}
	defer file.Close()

	buf := bufio.NewScanner(file)

	for buf.Scan() {
		m := metric.Metrics{}
		if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
			return fmt.Errorf("cannot decode metric: %w", err)
		}
		if err := s.Storage.SetMetric(m); err != nil {
			return fmt.Errorf("cannot set metric: %w", err)
		}
	}

	if buf.Err() != nil {
		return fmt.Errorf("cannot scan file: %w", buf.Err())
	}

	return nil
}

func (s *Server) Dumpfile() error {
	dumpFileSync.Lock()
	defer dumpFileSync.Unlock()

	file, err := os.OpenFile(s.Cfg.FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("cannot open/create dump file for write: %w", err)
	}
	defer file.Close()

	buf := bufio.NewWriter(file)
	defer buf.Flush()

	encoder := json.NewEncoder(buf)

	metrics := s.Storage.GetAllMetrics()
	for _, v := range metrics {
		if err := encoder.Encode(&v); err != nil {
			return fmt.Errorf("cannot encode metric: %w", err)
		}
	}

	return nil
}

func (s *Server) Syncfile() {
	syncTicker := time.NewTicker(s.Cfg.GetStoreInterval())
	defer syncTicker.Stop()

	for range syncTicker.C {
		if err := s.Dumpfile(); err != nil {
			log.Fatalf("cannot dump file: %s", err)
		}
	}
}

func NewDumper() {
}
