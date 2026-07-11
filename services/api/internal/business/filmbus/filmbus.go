// Package filmbus provides business logic related to films.
package filmbus

import (
	"context"
	"errors"

	"github.com/andrew-hayworth22/critiquefi-service/internal/business"
	"github.com/andrew-hayworth22/critiquefi-service/internal/models"
	"github.com/andrew-hayworth22/critiquefi-service/internal/store"
)

// Store defines the storage logic for films.
type Store interface {
	CreateFilm(ctx context.Context, film models.NewFilm) (id int64, err error)
	GetFilmByID(ctx context.Context, id int64) (models.Film, error)
}

type Bus struct {
	store Store
}

// New creates a new film business logic package.
func New(store Store) *Bus {
	return &Bus{store: store}
}

// CreateFilm creates a new film.
func (b *Bus) CreateFilm(ctx context.Context, film models.NewFilm) (int64, error) {
	// Validate new film
	if err := film.Validate(); err != nil {
		return 0, err
	}

	id, err := b.store.CreateFilm(ctx, film)
	if err != nil {
		if errors.Is(err, store.ErrDuplicate) {
			err = business.ErrDuplicate
		}
		return 0, err
	}

	return id, nil
}

// GetFilmByID returns a film by its ID.
func (b *Bus) GetFilmByID(ctx context.Context, id int64) (models.Film, error) {
	film, err := b.store.GetFilmByID(ctx, id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			err = business.ErrNotFound
		}
		return models.Film{}, err
	}
	return film, nil
}
