package server

import (
	"context"
	"net/http"
)

func (s *Server) PingDBHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.HandlerTimeout)
	defer cancel()

	w.Header().Set("Content-Type", "application/json")

	if err := s.Storage.Ping(ctx); err != nil {
		sendError(w, "cannot ping DB", err)
		return
	}
}
