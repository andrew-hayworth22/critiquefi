package store

import "errors"

var (
	ErrInternal            = errors.New("internal error")
	ErrNotFound            = errors.New("not found")
	ErrDuplicate           = errors.New("duplicate record")
	ErrForeignKeyViolation = errors.New("foreign key violation")
	ErrNotNullViolation    = errors.New("not null violation")
)
