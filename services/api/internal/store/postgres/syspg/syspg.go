package syspg

import (
	"context"

	"github.com/andrew-hayworth22/critiquefi-service/internal/store"
	"github.com/jmoiron/sqlx"
)

// SysPg is the postgres implementation of the system store
type SysPg struct {
	db *sqlx.DB
}

// New creates a new system store
func New(db *sqlx.DB) *SysPg {
	return &SysPg{db: db}
}

// Ping tests the database connection
func (s *SysPg) Ping(ctx context.Context) error {
	if err := s.db.PingContext(ctx); err != nil {
		return store.ErrInternal
	}
	return nil
}
