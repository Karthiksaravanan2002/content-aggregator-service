package server

import (
	"context"
	"net/http"
	"time"

	"dev.azure.com/daimler-mic/content-aggregator/service/props"
	"go.uber.org/zap"
)

// use rocoli http client

type Server struct {
	cfg    *props.Config
	logger *zap.Logger
	mux    *http.ServeMux
	srv    *http.Server
}

func NewServer(cfg *props.Config, logger *zap.Logger) *Server {
	mux := http.NewServeMux()

	s := &Server{
		cfg:    cfg,
		logger: logger,
		mux:    mux,
		srv: &http.Server{
			Addr:         cfg.Server.Address,
			Handler:      mux,
			ReadTimeout:  time.Duration(cfg.Server.Timeout) * time.Second,
			WriteTimeout: time.Duration(cfg.Server.Timeout) * time.Second,
		},
	}

	return s
}

func (s *Server) Mux() *http.ServeMux {
	return s.mux
}

func (s *Server) Start() error {
	s.logger.Info("server starting",
		zap.String("address", s.cfg.Server.Address),
	)

	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown() {
	s.logger.Info("server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		s.logger.Error("server shutdown error", zap.Error(err))
	} else {
		s.logger.Info("server shutdown complete")
	}
}
