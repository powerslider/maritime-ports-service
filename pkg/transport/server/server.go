package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/powerslider/maritime-ports-service/pkg/configs"
)

// Server represents an HTTP server.
type Server struct {
	serverInst *http.Server
}

// NewServer constructs new HTTP server with the provided muxer.
func NewServer(
	config *configs.Config,
	muxer *mux.Router,
) *Server {
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler: muxer,
	}

	return &Server{
		serverInst: server,
	}
}

// Start starts the HTTP server.
func (s *Server) Start(ctx context.Context, errChan chan error) {
	log.Printf("[Start] HTTP serverInst is starting on %s:\n", s.serverInst.Addr)

	err := s.serverInst.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		errChan <- errors.WithStack(err)
	}
}

// Stop stops the HTTP server.
func (s *Server) Stop(ctx context.Context) error {
	log.Println("[Shutdown] HTTP serverInst is shutting down...")

	shutdownCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	err := s.serverInst.Shutdown(shutdownCtx)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Run manages the HTTP server lifecycle on start and on shutdown.
func (s *Server) Run(ctx context.Context) error {
	errChan := make(chan error)

	go s.Start(ctx, errChan)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigs:
		return s.Stop(ctx)
	case err := <-errChan:
		return errors.WithStack(err)
	}
}
