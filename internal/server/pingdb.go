package server

import (
	"fmt"
	"net/http"
)

func (s *Server) PingDBHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if s.db == nil {
		sendError(w, "cannot ping DB", fmt.Errorf("DB connection not open"))
		return
	}

	if err := s.pingDB(); err != nil {
		sendError(w, "cannot ping DB", err)
		return
	}
}
