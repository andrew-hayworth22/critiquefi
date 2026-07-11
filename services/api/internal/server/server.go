package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/andrew-hayworth22/critiquefi-service/internal/http/authhttp"
	"github.com/andrew-hayworth22/critiquefi-service/internal/http/filmhttp"
	"github.com/andrew-hayworth22/critiquefi-service/internal/http/syshttp"
	"github.com/andrew-hayworth22/critiquefi-service/internal/middleware"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

type Server struct {
	httpServer *http.Server
}

// Config represents the server's configuration
type Config struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// Dependencies represent all the packages needed to build the server's endpoints
type Dependencies struct {
	Logger *slog.Logger

	// Handler packages
	AuthHandler *authhttp.Handler
	SysHandler  *syshttp.Handler
	FilmHandler *filmhttp.Handler

	// Middleware packages
	AuthMiddleware *middleware.AuthMiddleware
}

func New(config Config, dependencies Dependencies) *Server {
	r := chi.NewRouter()
	r.Use(httprate.LimitByRealIP(100, time.Second))
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(middleware.RequestLogger(dependencies.Logger))
	r.Use(chiMiddleware.Recoverer)

	registerRoutes(r, dependencies)

	return &Server{
		httpServer: &http.Server{
			Addr:         ":" + config.Addr,
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
			IdleTimeout:  config.IdleTimeout,
			Handler:      r,
		},
	}
}

// Start begins listening for incoming requests
func (s *Server) Start() error {
	if err := s.httpServer.ListenAndServe(); err != nil {
		return fmt.Errorf("error starting server: %w", err)
	}
	return nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("error shutting down server: %w", err)
	}
	return nil
}
