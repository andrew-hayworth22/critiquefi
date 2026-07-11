package authpg

import (
	"context"
	"fmt"
	"time"

	"github.com/andrew-hayworth22/critiquefi-service/internal/models"
	"github.com/andrew-hayworth22/critiquefi-service/internal/store"
	"github.com/andrew-hayworth22/critiquefi-service/internal/store/postgres"
	"github.com/jmoiron/sqlx"
)

// AuthPg is the postgres implementation of the auth store
type AuthPg struct {
	db *sqlx.DB
}

// New creates a new auth postgres store
func New(db *sqlx.DB) *AuthPg {
	return &AuthPg{db: db}
}

// userRow represents a user row in the database
type userRow struct {
	ID           int64     `db:"id"`
	Email        string    `db:"email"`
	DisplayName  string    `db:"display_name"`
	Name         string    `db:"name"`
	PasswordHash string    `db:"password_hash"`
	IsAdmin      bool      `db:"is_admin"`
	LastLogin    time.Time `db:"last_login"`
	IsActive     bool      `db:"is_active"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func (usr userRow) toModel() models.User {
	return models.User{
		ID:           usr.ID,
		Email:        usr.Email,
		DisplayName:  usr.DisplayName,
		Name:         usr.Name,
		IsAdmin:      usr.IsAdmin,
		PasswordHash: usr.PasswordHash,
		IsActive:     usr.IsActive,
	}
}

// refreshTokenRow represents a refresh token row in the database
type refreshTokenRow struct {
	TokenHash  string     `db:"token_hash"`
	UserID     int64      `db:"user_id"`
	UserAgent  string     `db:"user_agent"`
	LastUsedAt *time.Time `db:"last_used_at"`
	ExpiresAt  time.Time  `db:"expires_at"`
	CreatedAt  time.Time  `db:"created_at"`
}

func (rt *refreshTokenRow) toModel() models.RefreshToken {
	return models.RefreshToken{
		TokenHash: rt.TokenHash,
		UserID:    rt.UserID,
		UserAgent: rt.UserAgent,
		ExpiresAt: rt.ExpiresAt,
		CreatedAt: rt.CreatedAt,
	}
}

// userFieldsTaken represents the result of a check for taken fields
type userFieldsTaken struct {
	EmailTaken       bool `db:"email_taken"`
	DisplayNameTaken bool `db:"display_name_taken"`
}

func (uft *userFieldsTaken) toModel() models.UserFieldsTaken {
	return models.UserFieldsTaken{
		EmailTaken:       uft.EmailTaken,
		DisplayNameTaken: uft.DisplayNameTaken,
	}
}

// passwordResetTokenRow represents a password reset token row in the database
type passwordResetTokenRow struct {
	TokenHash string    `db:"token_hash"`
	UserID    int64     `db:"user_id"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}

func (prt *passwordResetTokenRow) toModel() models.PasswordResetToken {
	return models.PasswordResetToken{
		TokenHash: prt.TokenHash,
		UserID:    prt.UserID,
		ExpiresAt: prt.ExpiresAt,
		CreatedAt: prt.CreatedAt,
	}
}

// CreateUser creates a user record and returns the new user's ID
func (a *AuthPg) CreateUser(ctx context.Context, user models.NewUser) (int64, error) {
	const q = `
		INSERT INTO users (
			email,
			display_name,
			name,
			password_hash
		) VALUES (:email, :display_name, :name, :password_hash)
		RETURNING id;
	`

	rows, err := a.db.NamedQueryContext(ctx, q, map[string]any{
		"email":         user.Email,
		"display_name":  user.DisplayName,
		"name":          user.Name,
		"password_hash": user.PasswordHash,
	})
	if err != nil {
		return 0, fmt.Errorf("error creating user: %w", postgres.MapError(err))
	}
	defer postgres.CloseRows(rows)
	if !rows.Next() {
		return 0, store.ErrNotFound
	}

	var id int64
	if err := rows.Scan(&id); err != nil {
		return 0, fmt.Errorf("error scanning user ID: %w", postgres.MapError(err))
	}
	return id, nil
}

// GetUserByID fetches a user by their ID
func (a *AuthPg) GetUserByID(ctx context.Context, id int64) (models.User, error) {
	const q = `
		SELECT
			id,
			email,
			display_name,
			name,
			password_hash,
			is_admin,
			is_active,
			created_at,
			updated_at,
			last_login
		FROM users
		WHERE id = :id
	`

	rows, err := a.db.NamedQueryContext(ctx, q, map[string]any{
		"id": id,
	})
	if err != nil {
		return models.User{}, store.ErrNotFound
	}
	defer postgres.CloseRows(rows)
	if !rows.Next() {
		return models.User{}, store.ErrNotFound
	}

	var row userRow
	if err := rows.StructScan(&row); err != nil {
		return models.User{}, postgres.MapError(err)
	}

	return row.toModel(), nil
}

// GetUserByEmail fetches a user by their email
func (a *AuthPg) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	const q = `
		SELECT
			id,
			email,
			display_name,
			name,
			password_hash,
			is_admin,
			is_active,
			created_at,
			updated_at,
			last_login
		FROM users
		WHERE email = :email
	`

	rows, err := a.db.NamedQueryContext(ctx, q, map[string]any{
		"email": email,
	})
	if err != nil {
		return models.User{}, postgres.MapError(err)
	}
	defer postgres.CloseRows(rows)
	if !rows.Next() {
		return models.User{}, store.ErrNotFound
	}

	var row userRow
	if err := rows.StructScan(&row); err != nil {
		return models.User{}, postgres.MapError(err)
	}

	return row.toModel(), nil
}

// CheckTakenUserFields checks if another user takes the given fields
func (a *AuthPg) CheckTakenUserFields(ctx context.Context, request models.NewUserRequest) (models.UserFieldsTaken, error) {
	const q = `
		SELECT
			EXISTS (SELECT 1 FROM users WHERE email = :email) AS email_taken,
			EXISTS (SELECT 1 FROM users WHERE display_name = :display_name) AS display_name_taken
	`

	rows, err := a.db.NamedQueryContext(ctx, q, map[string]any{
		"email":        request.Email,
		"display_name": request.DisplayName,
	})
	if err != nil {
		return models.UserFieldsTaken{}, postgres.MapError(err)
	}
	defer postgres.CloseRows(rows)

	if !rows.Next() {
		return models.UserFieldsTaken{}, store.ErrNotFound
	}

	var uft userFieldsTaken
	if err := rows.StructScan(&uft); err != nil {
		return models.UserFieldsTaken{}, postgres.MapError(err)
	}

	return uft.toModel(), nil
}

// SetUserLastLogin sets the last login time for a user to now
func (a *AuthPg) SetUserLastLogin(ctx context.Context, id int64) error {
	const q = `
		UPDATE users
		SET last_login = NOW()
		WHERE id = :id
	`

	_, err := a.db.NamedExecContext(ctx, q, map[string]any{
		"id": id,
	})
	return postgres.MapError(err)
}

// SetUserPassword sets a user's password
func (a *AuthPg) SetUserPassword(ctx context.Context, userID int64, passwordHash string) error {
	const q = `
		UPDATE users
		SET password_hash = :password_hash
		WHERE id = :id
	`

	_, err := a.db.NamedExecContext(ctx, q, map[string]any{
		"id":            userID,
		"password_hash": passwordHash,
	})
	return postgres.MapError(err)
}

// CreateRefreshToken creates a new refresh token
func (a *AuthPg) CreateRefreshToken(ctx context.Context, refreshToken models.RefreshToken) error {
	const q = `
		INSERT INTO refresh_tokens (
			token_hash,
			user_id,
			user_agent,
			expires_at
		) VALUES (:token_hash, :user_id, :user_agent, :expires_at);
	`

	_, err := a.db.NamedExecContext(ctx, q, map[string]any{
		"token_hash": refreshToken.TokenHash,
		"user_id":    refreshToken.UserID,
		"user_agent": refreshToken.UserAgent,
		"expires_at": refreshToken.ExpiresAt,
	})
	return postgres.MapError(err)
}

// GetRefreshToken fetches a refresh token by its hash
func (a *AuthPg) GetRefreshToken(ctx context.Context, tokenHash string) (models.RefreshToken, error) {
	const q = `
		SELECT
			token_hash,
			user_id,
			user_agent,
			expires_at,
			created_at
		FROM refresh_tokens
		WHERE token_hash = :token_hash
	`

	rows, err := a.db.NamedQueryContext(ctx, q, map[string]any{
		"token_hash": tokenHash,
	})
	if err != nil {
		return models.RefreshToken{}, postgres.MapError(err)
	}
	defer postgres.CloseRows(rows)
	if !rows.Next() {
		return models.RefreshToken{}, store.ErrNotFound
	}

	var row refreshTokenRow
	if err := rows.StructScan(&row); err != nil {
		return models.RefreshToken{}, postgres.MapError(err)
	}
	return row.toModel(), nil
}

// DeleteRefreshToken deletes a refresh token by its hash
func (a *AuthPg) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	const q = `
		DELETE FROM
			refresh_tokens
		WHERE token_hash = :token_hash
	`

	rows, err := a.db.NamedQueryContext(ctx, q, map[string]any{
		"token_hash": tokenHash,
	})
	if err != nil {
		return postgres.MapError(err)
	}
	defer postgres.CloseRows(rows)
	return nil
}

// CreatePasswordResetToken creates a new password reset token
func (a *AuthPg) CreatePasswordResetToken(ctx context.Context, token models.PasswordResetToken) error {
	const q = `
		INSERT INTO password_reset_tokens (
		    token_hash,
			user_id,
			expires_at
		) VALUES (:token_hash, :user_id, :expires_at);
	`
	_, err := a.db.NamedExecContext(ctx, q, map[string]any{
		"token_hash": token.TokenHash,
		"user_id":    token.UserID,
		"expires_at": token.ExpiresAt,
	})
	return postgres.MapError(err)
}

// GetPasswordResetToken fetches a reset password token by its hash
func (a *AuthPg) GetPasswordResetToken(ctx context.Context, tokenHash string) (models.PasswordResetToken, error) {
	const q = `
		SELECT
			token_hash,
			user_id,
			expires_at,
			created_at
		FROM password_reset_tokens
		WHERE token_hash = :token_hash
	`

	rows, err := a.db.NamedQueryContext(ctx, q, map[string]any{
		"token_hash": tokenHash,
	})
	if err != nil {
		return models.PasswordResetToken{}, postgres.MapError(err)
	}
	defer postgres.CloseRows(rows)
	if !rows.Next() {
		return models.PasswordResetToken{}, store.ErrNotFound
	}

	var row passwordResetTokenRow
	if err := rows.StructScan(&row); err != nil {
		return models.PasswordResetToken{}, postgres.MapError(err)
	}
	return row.toModel(), nil
}

// DeletePasswordResetTokensByUserID deletes all reset password tokens for a user
func (a *AuthPg) DeletePasswordResetTokensByUserID(ctx context.Context, userID int64) error {
	const q = `
		DELETE FROM
			password_reset_tokens
		WHERE user_id = :user_id
	`

	_, err := a.db.NamedExecContext(ctx, q, map[string]any{
		"user_id": userID,
	})
	return postgres.MapError(err)
}
