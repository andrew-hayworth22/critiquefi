// Package syshttp contains HTTP endpoints related to system checks
package syshttp

import (
	"context"
	"net/http"

	"github.com/andrew-hayworth22/critiquefi-service/pkg/httputil"
)

// Bus defines the business logic needed for system checks
type Bus interface {
	Ping(ctx context.Context) error
}

// Handler exposes HTTP endpoints related to system checks
type Handler struct {
	bus Bus
}

// New creates a new system HTTP handler
func New(bus Bus) *Handler {
	return &Handler{bus: bus}
}

// SystemCheckResponse represents the response of a system check
type SystemCheckResponse struct {
	Status string `json:"status"`
}

// Liveness checks the liveness of the HTTP server
func (h *Handler) Liveness(w http.ResponseWriter, r *http.Request) {
	httputil.WriteJSON(w, http.StatusOK, SystemCheckResponse{
		Status: "ok",
	})
}

// Readiness checks the readiness of the HTTP server and the database connection
func (h *Handler) Readiness(w http.ResponseWriter, r *http.Request) {
	if err := h.bus.Ping(r.Context()); err != nil {
		httputil.WriteServiceUnavailable(w)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, SystemCheckResponse{
		Status: "ok",
	})
}
