package business

import "errors"

var (
	ErrInternal           = errors.New("internal error")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrDuplicate          = errors.New("duplicate record")
	ErrNotFound           = errors.New("not found")
	ErrEmailSend          = errors.New("mail send error")
)
