package httpbase

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

const (
	defaultReadTimeout       = 15 * time.Second
	defaultReadHeaderTimeout = 15 * time.Second
	defaultWriteTimeout      = 15 * time.Second
	defaultIdleTimeout       = 60 * time.Second
)

type ConfigTimeout struct {
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

type Server struct {
	server *http.Server
}

func NewServer(mux http.Handler, cfg ConfigTimeout, opts ...serverOption) *Server {
	if cfg.ReadTimeout <= 0 {
		cfg.ReadTimeout = defaultReadTimeout
	}
	if cfg.ReadHeaderTimeout <= 0 {
		cfg.ReadHeaderTimeout = defaultReadHeaderTimeout
	}
	if cfg.WriteTimeout <= 0 {
		cfg.WriteTimeout = defaultWriteTimeout
	}
	if cfg.IdleTimeout <= 0 {
		cfg.IdleTimeout = defaultIdleTimeout
	}
	httpServer := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}
	s := &Server{
		server: httpServer,
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	slog.Info("Stop http server")
	return s.server.Shutdown(ctx)
}
