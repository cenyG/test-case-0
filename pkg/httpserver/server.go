package httpserver

import (
	"context"
	"net/http"
	"time"
)

const (
	_defaultAddr            = ":80"
	_defaultShutdownTimeout = 3 * time.Second
)

// Server -.
type Server struct {
	server *http.Server
	notify chan error
}

// New -.
func New(handler http.Handler, opts ...Option) *Server {
	httpServer := &http.Server{
		Handler: handler,
		Addr:    _defaultAddr,
	}

	s := &Server{
		server: httpServer,
		notify: make(chan error, 1),
	}

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) Start() {
	s.notify <- s.server.ListenAndServe()
	close(s.notify)
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), _defaultShutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
