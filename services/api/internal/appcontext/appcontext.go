// Package appcontext controls context values for the application.
package appcontext

import (
	"context"
	"log/slog"

	"github.com/andrew-hayworth22/critiquefi-service/internal/models"
)

// contextKey is a value for use with context.WithValue
type contextKey string

const (
	claimsKey contextKey = "claims"
	loggerKey contextKey = "logger"
)

// SetClaims sets the claims in the context.
func SetClaims(ctx context.Context, claims models.Claims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

// GetClaims returns the claims from the context.
func GetClaims(ctx context.Context) (models.Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(models.Claims)
	return claims, ok
}

// SetLogger sets the logger in the context.
func SetLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// GetLogger returns the logger from the context.
func GetLogger(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(loggerKey).(*slog.Logger)
	if !ok {
		return slog.Default()
	}
	return logger
}
