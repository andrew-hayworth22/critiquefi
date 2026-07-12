// Package authbus provides authbus-related business logic.
package authbus

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/andrew-hayworth22/critiquefi/services/api/internal/business"
	"github.com/andrew-hayworth22/critiquefi/services/api/internal/mail"
	"github.com/andrew-hayworth22/critiquefi/services/api/internal/models"
	"github.com/andrew-hayworth22/critiquefi/services/api/internal/store"
	"github.com/andrew-hayworth22/critiquefi/services/api/pkg/crypto"
	"github.com/golang-jwt/jwt/v5"
)

// Store defines the logic needed for authbus storage
type Store interface {
	CreateUser(ctx context.Context, user models.NewUser) (id int64, err error)
	GetUserByID(ctx context.Context, id int64) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (user models.User, err error)
	CheckTakenUserFields(ctx context.Context, newUserRequest models.NewUserRequest) (fields models.UserFieldsTaken, err error)
	SetUserLastLogin(ctx context.Context, id int64) error
	SetUserPassword(ctx context.Context, id int64, password string) error

	CreateRefreshToken(ctx context.Context, refreshToken models.RefreshToken) (err error)
	GetRefreshToken(ctx context.Context, tokenHash string) (models.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error

	CreatePasswordResetToken(ctx context.Context, token models.PasswordResetToken) error
	GetPasswordResetToken(ctx context.Context, tokenHash string) (models.PasswordResetToken, error)
	DeletePasswordResetTokensByUserID(ctx context.Context, userID int64) error
}

// Mailer defines the logic needed for sending auth-related emails
type Mailer interface {
	SendWelcome(ctx context.Context, recipient mail.Recipient) error
	SendPasswordReset(ctx context.Context, recipient mail.Recipient, token string) error
}

// jwtClaims defines the claims to be stored in the JWT
type jwtClaims struct {
	jwt.RegisteredClaims
	UserID  int64  `json:"user_id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
}

func (c jwtClaims) toModel() models.Claims {
	return models.Claims{
		UserID:  c.UserID,
		Email:   c.Email,
		IsAdmin: c.IsAdmin,
	}
}

// Bus defines the auth business logic
type Bus struct {
	logger                *slog.Logger
	store                 Store
	mailer                Mailer
	accessTokenKey        []byte
	accessTokenTTL        time.Duration
	refreshTokenTTL       time.Duration
	passwordResetTokenTTL time.Duration
}

// BusConfig defines the auth business logic configuration
type BusConfig struct {
	Logger                *slog.Logger
	Store                 Store
	Mailer                Mailer
	AccessTokenKey        string
	AccessTokenTTL        time.Duration
	RefreshTokenTTL       time.Duration
	PasswordResetTokenTTL time.Duration
}

// New creates a new auth business logic package
func New(cfg BusConfig) *Bus {
	return &Bus{
		logger:                cfg.Logger,
		store:                 cfg.Store,
		mailer:                cfg.Mailer,
		accessTokenKey:        []byte(cfg.AccessTokenKey),
		accessTokenTTL:        cfg.AccessTokenTTL,
		refreshTokenTTL:       cfg.RefreshTokenTTL,
		passwordResetTokenTTL: cfg.PasswordResetTokenTTL,
	}
}

// Register creates a user and starts an authenticated session
func (b *Bus) Register(ctx context.Context, newUserRequest models.NewUserRequest, userAgent string, remember bool) (accessToken string, refreshToken string, err error) {
	// Validate new user
	if err = newUserRequest.Validate(); err != nil {
		return
	}

	taken, err := b.store.CheckTakenUserFields(ctx, newUserRequest)
	if err != nil {
		return
	}
	ve := models.ValidationErrors{}
	if taken.EmailTaken {
		ve.Add("email", "email already taken")
	}
	if taken.DisplayNameTaken {
		ve.Add("display_name", "display name already taken")
	}
	if ve.Any() {
		err = ve
		return
	}

	// Hash password
	hashedPassword, err := crypto.HashPassword(newUserRequest.Password)
	if err != nil {
		return
	}

	newUser := models.NewUser{
		Email:        newUserRequest.Email,
		DisplayName:  newUserRequest.DisplayName,
		Name:         newUserRequest.Name,
		PasswordHash: hashedPassword,
	}

	// Create user
	id, err := b.store.CreateUser(ctx, newUser)
	if err != nil {
		if errors.Is(err, store.ErrDuplicate) {
			err = business.ErrDuplicate
		}
		return
	}

	// Fetch user and generate tokens
	user, err := b.store.GetUserByID(ctx, id)
	if err != nil {
		return
	}

	accessToken, err = b.GenerateAccessToken(user)
	if err != nil {
		return
	}

	if !remember {
		return
	}

	refreshToken, err = b.GenerateRefreshToken(ctx, user, userAgent)
	if err != nil {
		return
	}

	// Send a welcome email and gracefully handle errors
	err = b.mailer.SendWelcome(ctx, mail.Recipient{
		Name:    user.Name,
		Address: user.Email,
	})
	if err != nil {
		b.logger.Error("failed to send welcome email", "err", err)
		err = nil
	}

	return
}

// Login authenticates a user and returns an access token and refresh token
func (b *Bus) Login(ctx context.Context, email, password, userAgent string, remember bool) (accessToken string, refreshToken string, err error) {
	user, err := b.store.GetUserByEmail(ctx, email)
	if err != nil || !user.IsActive {
		err = business.ErrInvalidCredentials
		return
	}

	if crypto.CompareHash(user.PasswordHash, password) != nil {
		err = business.ErrInvalidCredentials
		return
	}

	accessToken, err = b.GenerateAccessToken(user)
	if err != nil {
		return
	}

	err = b.store.SetUserLastLogin(ctx, user.ID)
	if err != nil {
		return
	}

	if !remember {
		return
	}

	refreshToken, err = b.GenerateRefreshToken(ctx, user, userAgent)
	if err != nil {
		return
	}

	return
}

// Logout invalidates a refresh token
func (b *Bus) Logout(ctx context.Context, refreshToken string) error {
	refreshToken = crypto.HashToken(refreshToken)

	if err := b.store.DeleteRefreshToken(ctx, refreshToken); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil
		}
		return err
	}
	return nil
}

// Refresh refreshes an access token using a refresh token
func (b *Bus) Refresh(ctx context.Context, refreshToken string) (accessToken, newRefreshToken string, err error) {
	// Hash provided refresh token
	refreshToken = crypto.HashToken(refreshToken)

	// Fetch refresh token
	token, err := b.store.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		err = business.ErrInvalidToken
		return "", "", err
	}

	// Revoke provided refresh token
	if err := b.store.DeleteRefreshToken(ctx, refreshToken); err != nil {
		return "", "", err
	}

	// Check for expired token
	if time.Now().UTC().After(token.ExpiresAt) {
		err = business.ErrInvalidToken
		return
	}

	// Fetch the user associated with the refresh token
	user, err := b.store.GetUserByID(ctx, token.UserID)
	if err != nil || !user.IsActive {
		err = business.ErrInvalidToken
		return
	}

	// Generate a new access token
	accessToken, err = b.GenerateAccessToken(user)
	if err != nil {
		return
	}

	// Rotate refresh tokens
	newRefreshToken, err = b.GenerateRefreshToken(ctx, user, token.UserAgent)
	if err != nil {
		return
	}
	return
}

// GenerateAccessToken generates an access token for a user
func (b *Bus) GenerateAccessToken(user models.User) (string, error) {
	claims := &jwtClaims{
		UserID:  user.ID,
		Email:   user.Email,
		IsAdmin: user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(b.accessTokenTTL)),
			Issuer:    "critiquefi",
			Subject:   fmt.Sprint(user.ID),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(b.accessTokenKey)
}

// ValidateAccessToken parses and validates an access token, then returns the token claims
func (b *Bus) ValidateAccessToken(accessToken string) (models.Claims, error) {
	t, err := jwt.ParseWithClaims(accessToken, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return b.accessTokenKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return models.Claims{}, business.ErrInvalidToken
	}

	claims, ok := t.Claims.(*jwtClaims)
	if !ok || !t.Valid {
		return models.Claims{}, business.ErrInvalidToken
	}

	claimsModel := claims.toModel()
	return claimsModel, nil
}

// GenerateRefreshToken generates a refresh token for a user and sends an email
func (b *Bus) GenerateRefreshToken(ctx context.Context, user models.User, userAgent string) (string, error) {
	refreshToken, err := crypto.RandomString(32)
	if err != nil {
		return "", err
	}

	hashedRefreshToken := crypto.HashToken(refreshToken)

	token := models.RefreshToken{
		TokenHash: hashedRefreshToken,
		UserID:    user.ID,
		UserAgent: userAgent,
		ExpiresAt: time.Now().Add(b.refreshTokenTTL).UTC(),
		CreatedAt: time.Now().UTC(),
	}

	if b.store.CreateRefreshToken(ctx, token) != nil {
		return "", err
	}

	return refreshToken, nil
}

// ForgotPassword sends a password reset email to a user
func (b *Bus) ForgotPassword(ctx context.Context, email string) error {
	user, err := b.store.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return business.ErrNotFound
		}
		return business.ErrInternal
	}

	if err = b.store.DeletePasswordResetTokensByUserID(ctx, user.ID); err != nil {
		return business.ErrInternal
	}

	resetToken, err := crypto.RandomString(32)
	if err != nil {
		return err
	}
	hashedResetToken := crypto.HashToken(resetToken)
	fmt.Println("token=" + resetToken + ";hashed=" + hashedResetToken)

	token := models.PasswordResetToken{
		TokenHash: hashedResetToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(b.passwordResetTokenTTL).UTC(),
		CreatedAt: time.Now().UTC(),
	}

	if err := b.store.CreatePasswordResetToken(ctx, token); err != nil {
		return business.ErrInternal
	}

	if err := b.mailer.SendPasswordReset(ctx, mail.Recipient{
		Name:    user.Name,
		Address: user.Email,
	}, resetToken); err != nil {
		return err
	}

	return nil
}

// ResetPassword resets a user's password using a password reset token'
func (b *Bus) ResetPassword(ctx context.Context, token string, newPassword string) error {
	hashedToken := crypto.HashToken(token)
	fmt.Println("token=" + token + ";hashed=" + hashedToken)
	tok, err := b.store.GetPasswordResetToken(ctx, hashedToken)
	if err != nil {
		return business.ErrInvalidToken
	}
	if time.Now().UTC().After(tok.ExpiresAt) {
		return business.ErrInvalidToken
	}

	hashedPassword, err := crypto.HashPassword(newPassword)
	if err != nil {
		return business.ErrInternal
	}

	err = b.store.SetUserPassword(ctx, tok.UserID, hashedPassword)
	if err != nil {
		return business.ErrInternal
	}

	err = b.store.DeletePasswordResetTokensByUserID(ctx, tok.UserID)
	if err != nil {
		return business.ErrInternal
	}

	return nil
}
