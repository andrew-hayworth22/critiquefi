// Package authhttp provides auth-related HTTP handlers.
package authhttp

import (
	"context"
	"errors"
	"net/http"

	"github.com/andrew-hayworth22/critiquefi/services/api/internal/appcontext"
	"github.com/andrew-hayworth22/critiquefi/services/api/internal/business"
	"github.com/andrew-hayworth22/critiquefi/services/api/internal/models"
	"github.com/andrew-hayworth22/critiquefi/services/api/pkg/httputil"
)

// Bus defines the business logic needed for auth handlers
type Bus interface {
	Register(ctx context.Context, user models.NewUserRequest, userAgent string, remember bool) (string, string, error)
	Login(ctx context.Context, email, password string, userAgent string, remember bool) (string, string, error)
	Logout(ctx context.Context, refreshToken string) error
	Refresh(ctx context.Context, refreshToken string) (string, string, error)
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token string, newPassword string) error
}

// Handler exposes HTTP endpoints related to authentication
type Handler struct {
	bus                      Bus
	refreshTokenCookieName   string
	refreshTokenCookieDomain string
}

// New creates a new authbus HTTP handler
func New(bus Bus, refreshTokenCookieName, refreshTokenCookieDomain string) *Handler {
	return &Handler{
		bus:                      bus,
		refreshTokenCookieName:   refreshTokenCookieName,
		refreshTokenCookieDomain: refreshTokenCookieDomain,
	}
}

// RegisterRequest represents a request to register a new user
type RegisterRequest struct {
	Email           string `json:"email"`
	DisplayName     string `json:"display_name"`
	Name            string `json:"name"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Remember        bool   `json:"remember"`
}

// ToModel converts a RegisterRequest to a NewUserRequest
func (r *RegisterRequest) ToModel() models.NewUserRequest {
	return models.NewUserRequest{
		Email:           r.Email,
		DisplayName:     r.DisplayName,
		Name:            r.Name,
		Password:        r.Password,
		ConfirmPassword: r.ConfirmPassword,
	}
}

// LoginRequest represents a request to login a user
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

// ForgotPasswordRequest represents a request to reset a user's password
type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

// ResetPasswordRequest represents a request to reset a user's password
type ResetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

// AuthenticationResponse represents a response to an authentication request
type AuthenticationResponse struct {
	AccessToken string `json:"access_token"`
}

// Register handles a request to register a new user
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	logger := appcontext.GetLogger(r.Context())

	var req RegisterRequest
	if !httputil.DecodeRequest(w, r, &req) {
		return
	}

	accessToken, refreshToken, err := h.bus.Register(r.Context(), req.ToModel(), r.UserAgent(), req.Remember)
	if err != nil {
		switch {
		case errors.As(err, &models.ValidationErrors{}):
			httputil.WriteUnprocessable(w, err)
			return
		case errors.Is(err, business.ErrDuplicate):
			httputil.WriteConflict(w)
			return
		default:
			logger.Error("registration failed", "error", err)
			httputil.WriteInternalError(w)
			return
		}
	}

	if refreshToken != "" {
		h.setRefreshTokenCookie(w, refreshToken)
	}

	httputil.WriteJSON(w, http.StatusCreated, AuthenticationResponse{
		AccessToken: accessToken,
	})
}

// Login handles a request to login a user
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	logger := appcontext.GetLogger(r.Context())

	var req LoginRequest
	if !httputil.DecodeRequest(w, r, &req) {
		return
	}

	accessToken, refreshToken, err := h.bus.Login(r.Context(), req.Email, req.Password, r.UserAgent(), req.Remember)
	if err != nil {
		switch {
		case errors.Is(err, business.ErrInvalidCredentials):
			httputil.WriteUnauthorized(w)
			return
		default:
			logger.Error("login failed", "error", err)
			httputil.WriteInternalError(w)
			return
		}
	}

	if refreshToken != "" {
		h.setRefreshTokenCookie(w, refreshToken)
	}

	httputil.WriteJSON(w, http.StatusOK, AuthenticationResponse{
		AccessToken: accessToken,
	})
}

// Logout handles a request to logout a user
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	logger := appcontext.GetLogger(r.Context())

	cookie, err := r.Cookie(h.refreshTokenCookieName)
	if err != nil {
		httputil.WriteNoContent(w)
		return
	}

	if err := h.bus.Logout(r.Context(), cookie.Value); err != nil {
		logger.Error("logout failed", "error", err)
		httputil.WriteInternalError(w)
		return
	}

	h.clearRefreshTokenCookie(w)
	httputil.WriteNoContent(w)
}

// Refresh handles a request to refresh an access token
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(h.refreshTokenCookieName)
	if err != nil {
		httputil.WriteUnauthorized(w)
		return
	}

	accessToken, refreshToken, err := h.bus.Refresh(r.Context(), cookie.Value)
	if err != nil {
		httputil.WriteUnauthorized(w)
		return
	}

	if refreshToken != "" {
		h.setRefreshTokenCookie(w, refreshToken)
	}

	httputil.WriteJSON(w, http.StatusOK, AuthenticationResponse{
		AccessToken: accessToken,
	})
}

// ForgotPassword creates a new password reset token
func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	logger := appcontext.GetLogger(r.Context())

	var req ForgotPasswordRequest
	if !httputil.DecodeRequest(w, r, &req) {
		return
	}

	// Ignore the error, we don't want to communicate to the user if the email is valid or not
	err := h.bus.ForgotPassword(r.Context(), req.Email)
	if err != nil {
		logger.Error("forgot password failed", "error", err)
	}
	httputil.WriteNoContent(w)
}

// ResetPassword resets a user's password
func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if !httputil.DecodeRequest(w, r, &req) {
		return
	}

	err := h.bus.ResetPassword(r.Context(), req.Token, req.Password)
	if err != nil {
		if errors.Is(err, business.ErrInvalidToken) {
			httputil.WriteBadRequest(w, "invalid/expired token")
			return
		}
		httputil.WriteInternalError(w)
	}
	httputil.WriteNoContent(w)
}

// setRefreshTokenCookie sets the refresh token cookie
func (h *Handler) setRefreshTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     h.refreshTokenCookieName,
		Value:    token,
		Domain:   h.refreshTokenCookieDomain,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

// clearRefreshTokenCookie clears the refresh token cookie
func (h *Handler) clearRefreshTokenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     h.refreshTokenCookieName,
		Domain:   h.refreshTokenCookieDomain,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})
}
