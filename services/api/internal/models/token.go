package models

import "time"

// Claims represent the claims in the application's JWT
type Claims struct {
	UserID  int64  `json:"uid"`
	Email   string `json:"mail"`
	IsAdmin bool   `json:"is_admin"`
}

// RefreshToken represents a token used to refresh the user's access token
type RefreshToken struct {
	TokenHash string
	UserID    int64
	UserAgent string
	ExpiresAt time.Time
	CreatedAt time.Time
}

// PasswordResetToken represents a token used to reset the user's password
type PasswordResetToken struct {
	TokenHash string
	UserID    int64
	ExpiresAt time.Time
	CreatedAt time.Time
}
