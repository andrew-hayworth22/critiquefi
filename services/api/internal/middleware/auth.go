package middleware

import (
	"fmt"
	"net/http"

	"github.com/andrew-hayworth22/critiquefi-service/internal/appcontext"
	"github.com/andrew-hayworth22/critiquefi-service/internal/models"
	"github.com/andrew-hayworth22/critiquefi-service/pkg/httputil"
)

// Bus defines the business logic needed for authbus middleware
type Bus interface {
	ValidateAccessToken(accessToken string) (models.Claims, error)
}

// AuthMiddleware provides middleware for authenticating/authorizing requests
type AuthMiddleware struct {
	bus Bus
}

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware(b Bus) *AuthMiddleware {
	return &AuthMiddleware{bus: b}
}

// Authenticate optionally authenticates the request and stores claims in the context
func (a *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := extractBearerToken(r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		claims, err := a.bus.ValidateAccessToken(token)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := appcontext.SetClaims(r.Context(), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ForceAuthentication ensures that the request is authenticated with a valid access token
func (a *AuthMiddleware) ForceAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := appcontext.GetClaims(r.Context())
		if !ok {
			httputil.WriteUnauthorized(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ForceAdmin ensures that the request is authenticated with an admin
func (a *AuthMiddleware) ForceAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := appcontext.GetClaims(r.Context())
		if !ok {
			httputil.WriteUnauthorized(w)
			return
		}
		if !claims.IsAdmin {
			httputil.WriteForbidden(w)
			return
		}
	})
}

// extractBearerToken extracts the access token from the authorization request header
func extractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing authorization header")
	}

	const prefix = "Bearer "
	if len(authHeader) < len(prefix) || authHeader[:len(prefix)] != prefix {
		return "", fmt.Errorf("invalid authorization header")
	}

	return authHeader[len(prefix):], nil
}
