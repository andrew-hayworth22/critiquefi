package server

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
)

// registerRoutes registers all routes for the application
func registerRoutes(r *chi.Mux, dependencies Dependencies) {
	r.Group(unauthenticatedRoutes(dependencies))
	r.Group(publicRoutes(dependencies))
	r.Group(protectedRoutes(dependencies))
	r.Group(adminRoutes(dependencies))
}

// unauthenticatedRoutes defines routes that will not attempt authentication
func unauthenticatedRoutes(dependencies Dependencies) func(chi.Router) {
	return func(r chi.Router) {
		// System Checks
		r.Get("/liveness", dependencies.SysHandler.Liveness)
		r.Get("/readiness", dependencies.SysHandler.Readiness)

		// Auth
		r.Post("/auth/register", dependencies.AuthHandler.Register)
		r.Post("/auth/login", dependencies.AuthHandler.Login)
		r.Post("/auth/refresh", dependencies.AuthHandler.Refresh)
		r.With(httprate.LimitByRealIP(1, time.Minute)).Post("/auth/forgot-password", dependencies.AuthHandler.ForgotPassword)
		r.With(httprate.LimitByRealIP(1, time.Minute)).Post("/auth/reset-password", dependencies.AuthHandler.ResetPassword)
	}
}

// publicRoutes defines routes that do not require authentication
func publicRoutes(dependencies Dependencies) func(chi.Router) {
	return func(r chi.Router) {
		r.Use(dependencies.AuthMiddleware.Authenticate)
	}
}

// protectedRoutes defines routes that require authentication
func protectedRoutes(dependencies Dependencies) func(chi.Router) {
	return func(r chi.Router) {
		r.Use(dependencies.AuthMiddleware.Authenticate)
		r.Use(dependencies.AuthMiddleware.ForceAuthentication)

		// Auth
		r.Post("/auth/logout", dependencies.AuthHandler.Logout)

		// Film
		r.Post("/films", dependencies.FilmHandler.CreateFilm)
		r.Get("/films/{id}", dependencies.FilmHandler.GetFilmById)
	}
}

// adminRoutes defines routes that
func adminRoutes(dependencies Dependencies) func(chi.Router) {
	return func(r chi.Router) {
		r.Use(dependencies.AuthMiddleware.Authenticate)
		r.Use(dependencies.AuthMiddleware.ForceAdmin)
	}
}
