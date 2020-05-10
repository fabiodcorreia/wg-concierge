package web

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fabiodcorreia/wg-concierge/internal/conf"
)

// HTTPServer is the application http server, it contains all the services needed as dependencies
type HTTPServer struct {
	server http.Server
	logger *log.Logger
}

// NewHTTPServer creates a new application server
func NewHTTPServer(cfg conf.App, logger *log.Logger, handler http.Handler) *HTTPServer {
	appServer := HTTPServer{
		server: http.Server{
			Addr:              cfg.Web.Host,              // IP:Port or Hostname:Port
			Handler:           handler,                   // Root HTTP Router
			ReadTimeout:       cfg.Web.ReadTimeout,       // the maximum duration for reading the entire request, including the body
			WriteTimeout:      cfg.Web.WriteTimeout,      // the maximum duration before timing out writes of the response
			IdleTimeout:       cfg.Web.IdleTimeout,       // the maximum amount of time to wait for the next request when keep-alive is enabled
			ReadHeaderTimeout: cfg.Web.ReadHeaderTimeout, // the amount of time allowed to read request headers
		},
		logger: logger,
	}
	return &appServer
}

// StartAndWait starts the application server and waits
func (h *HTTPServer) StartAndWait(shutdownTimeout time.Duration) error {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	serverErrors := make(chan error, 1)

	go func() {
		h.logger.Printf("Server listening on %s", h.server.Addr)
		serverErrors <- h.server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		h.logger.Printf("%v : Start shutdown", sig)

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		err := h.server.Shutdown(ctx)
		if err != nil {
			h.logger.Printf("Graceful shutdown did not complete in %v : %v", shutdownTimeout, err)
			err = h.server.Close()
		}

		switch {
		case sig == syscall.SIGSTOP:
			return errors.New("integrity issue caused shutdown")
		case err != nil:
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}
	return nil
}
