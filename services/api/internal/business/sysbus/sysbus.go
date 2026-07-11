// Package sysbus provides business logic for system-related operations.
package sysbus

import (
	"context"

	"github.com/andrew-hayworth22/critiquefi-service/internal/business"
)

// Store defines logic for system-related storage operations.
type Store interface {
	Ping(ctx context.Context) error
}

// Bus defines the system-related business logic.
type Bus struct {
	store Store
}

// New creates a new system-related business logic package.
func New(store Store) *Bus {
	return &Bus{store: store}
}

// Ping tests the database connection.
func (b *Bus) Ping(ctx context.Context) error {
	if err := b.store.Ping(ctx); err != nil {
		return business.ErrInternal
	}
	return nil
}
