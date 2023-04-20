package server

import "github.com/av-baran/ymetrics/internal/repository"

type Server struct {
	Storage repository.Storager
}

func New(s repository.Storager) *Server {
	return &Server{Storage: s}
}
