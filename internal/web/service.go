package web

import (
	"log"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Service
type Service struct {
	*chi.Mux
	logger *log.Logger
}

// NewService
func NewService(logger *log.Logger) *Service {
	s := Service{
		Mux:    chi.NewRouter(),
		logger: logger,
	}

	// Set generic middleware
	s.Mux.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: logger, NoColor: false}))
	s.Mux.Use(middleware.Recoverer)
	s.Mux.Use(middleware.RealIP)
	s.Mux.Use(middleware.StripSlashes)

	return &s
}
