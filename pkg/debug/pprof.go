// Package debug implements the pprof server.
package debug

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"
)

const pprofEndpoint = "/debug/pprof/"

type Server struct {
	server *http.Server
	port   int
}

func NewServer(port int) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	// Register pprof endpoints - they are automatically registered to the default ServeMux.
	mux.Handle(pprofEndpoint, http.DefaultServeMux)

	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- s.server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		return s.Shutdown(context.Background())
	case err := <-errChan:
		return err
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	return s.server.Shutdown(ctx)
}
