// Package filmhttp provides film-related HTTP handlers.
package filmhttp

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/andrew-hayworth22/critiquefi-service/internal/appcontext"
	"github.com/andrew-hayworth22/critiquefi-service/internal/business"
	"github.com/andrew-hayworth22/critiquefi-service/internal/models"
	"github.com/andrew-hayworth22/critiquefi-service/pkg/httputil"
	"github.com/go-chi/chi/v5"
)

// Bus defines the business logic needed for film-related HTTP handlers.
type Bus interface {
	CreateFilm(ctx context.Context, film models.NewFilm) (int64, error)
	GetFilmByID(ctx context.Context, id int64) (models.Film, error)
}

// Handler exposes HTTP endpoints related to films.
type Handler struct {
	bus Bus
}

// New creates a new film HTTP handler.
func New(bus Bus) *Handler {
	return &Handler{bus: bus}
}

// ExternalReference represents a reference to an external resource.
type ExternalReference struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// NewFilmRequest represents a request to create a new film.
type NewFilmRequest struct {
	FilmType           string              `json:"film_type"`
	Title              string              `json:"title"`
	Description        *string             `json:"description,omitempty"`
	ReleaseDate        string              `json:"release_date"`
	RuntimeMinutes     *int                `json:"runtime_minutes,omitempty"`
	ExternalReferences []ExternalReference `json:"external_references,omitempty"`
}

func (r *NewFilmRequest) toModel() (models.NewFilm, error) {
	ve := models.ValidationErrors{}

	filmType, err := models.ParseFilmType(r.FilmType)
	if err != nil {
		ve.Add("film_type", "invalid film type")
	}
	releaseDate, err := time.Parse("2006-01-02", r.ReleaseDate)
	if err != nil {
		ve.Add("release_date", "invalid release date")
	}
	if ve.Any() {
		return models.NewFilm{}, ve
	}

	var externalRefs []models.ExternalReference
	for _, ref := range r.ExternalReferences {
		externalRefs = append(externalRefs, models.ExternalReference{
			Name: ref.Name,
			URL:  ref.URL,
		})
	}

	return models.NewFilm{
		FilmType:           filmType,
		Title:              r.Title,
		Description:        r.Description,
		ReleaseDate:        releaseDate,
		RuntimeMinutes:     r.RuntimeMinutes,
		ExternalReferences: externalRefs,
	}, nil
}

type FilmResponse struct {
	ID                 int64               `json:"id"`
	FilmType           string              `json:"film_type"`
	Title              string              `json:"title"`
	Description        *string             `json:"description,omitempty"`
	ReleaseDate        time.Time           `json:"release_date"`
	RuntimeMinutes     *int                `json:"runtime_minutes,omitempty"`
	ExternalReferences []ExternalReference `json:"external_references,omitempty"`
}

func toResponse(f models.Film) FilmResponse {
	var externalRefs []ExternalReference
	for _, ref := range f.ExternalReferences {
		externalRefs = append(externalRefs, ExternalReference{
			Name: ref.Name,
			URL:  ref.URL,
		})
	}

	return FilmResponse{
		ID:                 f.ID,
		FilmType:           f.FilmType.String(),
		Title:              f.Title,
		Description:        f.Description,
		ReleaseDate:        f.ReleaseDate,
		RuntimeMinutes:     f.RuntimeMinutes,
		ExternalReferences: externalRefs,
	}
}

func (h *Handler) CreateFilm(w http.ResponseWriter, r *http.Request) {
	logger := appcontext.GetLogger(r.Context())

	var req NewFilmRequest
	if !httputil.DecodeRequest(w, r, &req) {
		return
	}

	newFilm, err := req.toModel()
	if err != nil {
		httputil.WriteUnprocessable(w, err)
		return
	}

	id, err := h.bus.CreateFilm(r.Context(), newFilm)
	if err != nil {
		switch {
		case errors.As(err, &models.ValidationErrors{}):
			httputil.WriteUnprocessable(w, err)
			return
		case errors.Is(err, business.ErrDuplicate):
			httputil.WriteConflict(w)
			return
		case errors.Is(err, business.ErrNotFound):
			httputil.WriteNotFound(w)
			return
		default:
			logger.Error("failed to create film", "error", err)
			httputil.WriteInternalError(w)
		}
	}

	httputil.WriteJSON(w, http.StatusCreated, map[string]int64{
		"id": id,
	})
}

func (h *Handler) GetFilmById(w http.ResponseWriter, r *http.Request) {
	logger := appcontext.GetLogger(r.Context())

	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		httputil.WriteNotFound(w)
		return
	}

	film, err := h.bus.GetFilmByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, business.ErrNotFound):
			httputil.WriteNotFound(w)
			return
		}
		logger.Error("failed to get film", "error", err)
		httputil.WriteInternalError(w)
	}

	httputil.WriteJSON(w, http.StatusOK, toResponse(film))
}
