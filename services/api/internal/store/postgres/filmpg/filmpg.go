package filmpg

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/andrew-hayworth22/critiquefi-service/internal/models"
	"github.com/andrew-hayworth22/critiquefi-service/internal/store"
	"github.com/andrew-hayworth22/critiquefi-service/internal/store/postgres"
	"github.com/jmoiron/sqlx"
)

// FilmPg is the postgres implementation of the film store
type FilmPg struct {
	db *sqlx.DB
}

// New creates a new film postgres store
func New(db *sqlx.DB) *FilmPg {
	return &FilmPg{db: db}
}

// filmRow represents a film row in the database
type filmRow struct {
	ID                 int64     `db:"id"`
	FilmType           string    `db:"film_type"`
	Title              string    `db:"title"`
	Description        *string   `db:"description"`
	ReleaseDate        time.Time `db:"release_date"`
	RuntimeMinutes     *int      `db:"runtime_minutes"`
	ExternalReferences []byte    `db:"external_references"`
	CreatedAt          time.Time `db:"created_at"`
	CreatedBy          *int64    `db:"created_by"`
	UpdatedAt          time.Time `db:"updated_at"`
	UpdatedBy          *int64    `db:"updated_by"`
}

func (f filmRow) toModel() (models.Film, error) {
	filmType, err := models.ParseFilmType(f.FilmType)
	if err != nil {
		return models.Film{}, fmt.Errorf("error parsing film type: %w", err)
	}

	var externalRefs []models.ExternalReference
	if err := json.Unmarshal(f.ExternalReferences, &externalRefs); err != nil {
		return models.Film{}, fmt.Errorf("error unmarshalling external references: %w", err)
	}

	return models.Film{
		ID:                 f.ID,
		FilmType:           filmType,
		Title:              f.Title,
		Description:        f.Description,
		ReleaseDate:        f.ReleaseDate,
		RuntimeMinutes:     f.RuntimeMinutes,
		ExternalReferences: externalRefs,
	}, nil
}

// CreateFilm creates a film record and returns that new film's ID
func (f *FilmPg) CreateFilm(ctx context.Context, newFilm models.NewFilm) (int64, error) {
	q := `
		INSERT INTO films (
			film_type,
			title,
		    description,
			release_date,
			runtime_minutes,
			external_references
		) VALUES (:film_type, :title, :description, :release_date, :runtime_minutes, :external_references)
		RETURNING id;
	`

	rows, err := f.db.NamedQueryContext(ctx, q, map[string]any{
		"film_type":           strings.ToUpper(newFilm.FilmType.String()),
		"title":               newFilm.Title,
		"description":         newFilm.Description,
		"release_date":        newFilm.ReleaseDate,
		"runtime_minutes":     newFilm.RuntimeMinutes,
		"external_references": newFilm.ExternalReferences,
	})
	if err != nil {
		return 0, fmt.Errorf("error creating film: %w", err)
	}
	defer postgres.CloseRows(rows)
	if !rows.Next() {
		return 0, store.ErrNotFound
	}

	var id int64
	if err := rows.Scan(&id); err != nil {
		return 0, fmt.Errorf("error scanning film ID: %w", err)
	}
	return id, nil
}

// GetFilmByID fetches a film by its ID
func (f *FilmPg) GetFilmByID(ctx context.Context, id int64) (models.Film, error) {
	q := `
		SELECT
			id,
			title,
			film_type,
			title,
			description,
			release_date,
			runtime_minutes,
			external_references,
			created_at,
			created_by,
			updated_at,
			updated_by
		FROM films
		WHERE id = :id
	`

	rows, err := f.db.NamedQueryContext(ctx, q, map[string]any{
		"id": id,
	})
	if err != nil {
		return models.Film{}, store.ErrNotFound
	}
	defer postgres.CloseRows(rows)
	if !rows.Next() {
		return models.Film{}, store.ErrNotFound
	}

	var row filmRow
	if err := rows.StructScan(&row); err != nil {
		return models.Film{}, postgres.MapError(err)
	}
	return row.toModel()
}
