package server

import (
	"net/http"
)

func (s *Server) PingDBHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := s.Storage.Ping(); err != nil {
		sendError(w, "cannot ping DB", err)
		return
	}
}
