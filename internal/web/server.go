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
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// HTTPServer is the application http server, it contains all the services
type HTTPServer struct {
	server http.Server
	logger *log.Logger
	router *chi.Mux
}

// NewHTTPServer creates a new application server with the provided configurations and root router.
func NewHTTPServer(cfg conf.App, logger *log.Logger) *HTTPServer {
	r := newRootRouter(logger)
	appServer := HTTPServer{
		server: http.Server{
			Addr:              cfg.Web.Host,              // IP:Port or Hostname:Port
			Handler:           r,                         // Root HTTP Router
			ReadTimeout:       cfg.Web.ReadTimeout,       // the maximum duration for reading the entire request, including the body
			WriteTimeout:      cfg.Web.WriteTimeout,      // the maximum duration before timing out writes of the response
			IdleTimeout:       cfg.Web.IdleTimeout,       // the maximum amount of time to wait for the next request when keep-alive is enabled
			ReadHeaderTimeout: cfg.Web.ReadHeaderTimeout, // the amount of time allowed to read request headers
		},
		logger: logger,
	}
	// For redirect from root to invitation form
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/invite", http.StatusPermanentRedirect)
	})
	appServer.router = r
	return &appServer
}

// AddService attaches a new Service to the root router
func (h *HTTPServer) AddService(path string, handler http.Handler) {
	h.logger.Printf("Route %s ready\n", path)
	h.router.Mount(path, handler)
}

// StartAndWait starts the application server and waits for a gracefull shutdown.
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

// newRootRouter init the root router with the default middleware for all the routes.
func newRootRouter(logger *log.Logger) *chi.Mux {
	root := chi.NewRouter()
	root.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: logger, NoColor: false}))
	root.Use(middleware.Recoverer)
	root.Use(middleware.RealIP)
	root.Use(middleware.StripSlashes)
	return root
}
