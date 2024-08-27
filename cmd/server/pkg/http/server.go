package http

import (
	"context"
	"net/http"
	"time"
)

const (
	defaultReadTimeout       = 15 * time.Second
	defaultReadHeaderTimeout = 15 * time.Second
	defaultWriteTimeout      = 15 * time.Second
	defaultIdleTimeout       = 60 * time.Second
)

type Server struct {
	server *http.Server
}

func NewServer(mux http.Handler, opts ...serverOption) *Server {
	httpServer := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       120 * time.Second,
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
	return s.server.Shutdown(ctx)
}
