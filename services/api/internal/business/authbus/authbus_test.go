package authbus_test

import (
	"context"
	"testing"
	"time"

	"github.com/andrew-hayworth22/critiquefi/services/api/internal/business"
	"github.com/andrew-hayworth22/critiquefi/services/api/internal/business/authbus"
	"github.com/andrew-hayworth22/critiquefi/services/api/internal/models"
	"github.com/andrew-hayworth22/critiquefi/services/api/internal/store"
	"github.com/andrew-hayworth22/critiquefi/services/api/internal/testutil"
	"github.com/andrew-hayworth22/critiquefi/services/api/pkg/crypto"
)

func TestBus_Register(t *testing.T) {
	tc := []struct {
		name                 string
		request              models.NewUserRequest
		userAgent            string
		remember             bool
		mockSetup            func(s *mockStore, m *mockMailer)
		expectedErr          error
		accessTokenExpected  bool
		refreshTokenExpected bool
	}{
		{
			name: "success: registered (don't remember)",
			request: models.NewUserRequest{
				Email:           "test.user@critiquefi.com",
				DisplayName:     "test.user",
				Name:            "Test User",
				Password:        "password",
				ConfirmPassword: "password",
			},
			userAgent: "test",
			remember:  false,
			mockSetup: func(s *mockStore, m *mockMailer) {
				s.On(CheckTakenUserFields, models.UserFieldsTaken{}, nil).
					On(CreateUser, int64(1), nil).
					On(GetUserByID, models.User{
						ID:           1,
						Email:        "user@critiquefi.com",
						DisplayName:  "user.name",
						Name:         "Test User",
						IsAdmin:      false,
						PasswordHash: "password",
						IsActive:     true,
					}, nil)
				m.On(SendWelcome, nil)
			},
			expectedErr:          nil,
			accessTokenExpected:  true,
			refreshTokenExpected: false,
		},
		{
			name: "success: registered (remember)",
			request: models.NewUserRequest{
				Email:           "test.user@critiquefi.com",
				DisplayName:     "test.user",
				Name:            "Test User",
				Password:        "password",
				ConfirmPassword: "password",
			},
			userAgent: "test",
			remember:  true,
			mockSetup: func(s *mockStore, m *mockMailer) {
				s.On(CheckTakenUserFields, models.UserFieldsTaken{}, nil).
					On(CreateUser, int64(1), nil).
					On(GetUserByID, models.User{
						ID:           1,
						Email:        "user@critiquefi.com",
						DisplayName:  "user.name",
						Name:         "Test User",
						IsAdmin:      false,
						PasswordHash: "password",
						IsActive:     true,
					}, nil).
					On(CreateRefreshToken, nil)
				m.On(SendWelcome, nil)
			},
			expectedErr:          nil,
			accessTokenExpected:  true,
			refreshTokenExpected: true,
		},
		{
			name: "error: validation errors - model validation",
			request: models.NewUserRequest{
				Email:           "critiquefi.com",
				DisplayName:     "er",
				Name:            "r",
				Password:        "ssword",
				ConfirmPassword: "password",
			},
			userAgent: "test",
			remember:  true,
			mockSetup: func(s *mockStore, m *mockMailer) {},
			expectedErr: models.ValidationErrors{
				"email":            "invalid email address",
				"display_name":     "display name must be between 3 and 50 characters long",
				"name":             "name must be between 3 and 50 characters long",
				"password":         "password must be between 8 and 64 characters long",
				"confirm_password": "passwords do not match",
			},
			accessTokenExpected:  false,
			refreshTokenExpected: false,
		},
		{
			name: "error: validation errors - duplicate fields",
			request: models.NewUserRequest{
				Email:           "test.user@critiquefi.com",
				DisplayName:     "test.user",
				Name:            "Test User",
				Password:        "password",
				ConfirmPassword: "password",
			},
			userAgent: "test",
			remember:  true,
			mockSetup: func(s *mockStore, m *mockMailer) {
				s.On(CheckTakenUserFields, models.UserFieldsTaken{
					DisplayNameTaken: true,
					EmailTaken:       true,
				}, nil)
			},
			expectedErr: models.ValidationErrors{
				"email":        "email already taken",
				"display_name": "display name already taken",
			},
			accessTokenExpected:  false,
			refreshTokenExpected: false,
		},
		{
			name: "error: duplicate creation",
			request: models.NewUserRequest{
				Email:           "test.user@critiquefi.com",
				DisplayName:     "test.user",
				Name:            "Test User",
				Password:        "password",
				ConfirmPassword: "password",
			},
			userAgent: "test",
			remember:  true,
			mockSetup: func(s *mockStore, m *mockMailer) {
				s.On(CheckTakenUserFields, models.UserFieldsTaken{}, nil).
					On(CreateUser, int64(0), store.ErrDuplicate)
			},
			expectedErr:          business.ErrDuplicate,
			accessTokenExpected:  false,
			refreshTokenExpected: false,
		},
	}

	for _, tc := range tc {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := newMockStore(t)
			m := newMockMailer(t)
			tc.mockSetup(s, m)

			svc := authbus.New(authbus.BusConfig{
				Store:           s,
				Mailer:          m,
				AccessTokenKey:  "key",
				AccessTokenTTL:  time.Hour,
				RefreshTokenTTL: time.Hour,
			})

			accessToken, refreshToken, err := svc.Register(context.Background(), tc.request, tc.userAgent, tc.remember)
			testutil.CheckErr(err, tc.expectedErr, t)

			if tc.accessTokenExpected && accessToken == "" {
				t.Error("expected access token, got nothing")
			}
			if !tc.accessTokenExpected && accessToken != "" {
				t.Errorf("expected no access token, got %v", accessToken)
			}
			if tc.refreshTokenExpected && refreshToken == "" {
				t.Error("expected refresh token, got nothing")
			}
			if !tc.refreshTokenExpected && refreshToken != "" {
				t.Errorf("expected no refresh token, got %v", refreshToken)
			}
		})
	}
}

func TestBus_Login(t *testing.T) {
	t.Parallel()

	hashedPassword, err := crypto.HashPassword("password")
	if err != nil {
		t.Fatal(err)
	}

	tc := []struct {
		name                 string
		email                string
		password             string
		userAgent            string
		remember             bool
		storeSetup           func(s *mockStore)
		expectedErr          error
		accessTokenExpected  bool
		refreshTokenExpected bool
	}{
		{
			name:      "success: logged in (don't remember)",
			email:     "user@critiquefi.com",
			password:  "password",
			userAgent: "test",
			remember:  false,
			storeSetup: func(s *mockStore) {
				s.On(GetUserByEmail, models.User{
					ID:           1,
					Email:        "user@critiquefi.com",
					DisplayName:  "test.user",
					Name:         "Test User",
					IsAdmin:      false,
					PasswordHash: hashedPassword,
					IsActive:     true,
				}, nil).On(SetUserLastLogin, nil)
			},
			expectedErr:          nil,
			accessTokenExpected:  true,
			refreshTokenExpected: false,
		},
		{
			name:      "success: logged in (remember)",
			email:     "user@critiquefi.com",
			password:  "password",
			userAgent: "test",
			remember:  true,
			storeSetup: func(s *mockStore) {
				s.On(GetUserByEmail, models.User{
					ID:           1,
					Email:        "user@critiquefi.com",
					DisplayName:  "test.user",
					Name:         "Test User",
					IsAdmin:      false,
					PasswordHash: hashedPassword,
					IsActive:     true,
				}, nil).On(SetUserLastLogin, nil).On(CreateRefreshToken, nil)
			},
			expectedErr:          nil,
			accessTokenExpected:  true,
			refreshTokenExpected: true,
		},
		{
			name:      "error: wrong mail",
			email:     "nonexistent@critiquefi.com",
			password:  "password",
			userAgent: "test",
			remember:  false,
			storeSetup: func(s *mockStore) {
				s.On(GetUserByEmail, models.User{}, store.ErrNotFound)
			},
			expectedErr:          business.ErrInvalidCredentials,
			accessTokenExpected:  false,
			refreshTokenExpected: false,
		},
		{
			name:      "error: wrong password",
			email:     "user@critiquefi.com",
			password:  "password",
			userAgent: "test",
			remember:  true,
			storeSetup: func(s *mockStore) {
				s.On(GetUserByEmail, models.User{
					ID:           1,
					Email:        "user@critiquefi.com",
					DisplayName:  "test.user",
					Name:         "Test User",
					IsAdmin:      false,
					PasswordHash: "wrongpassword",
					IsActive:     true,
				}, nil)
			},
			expectedErr:          business.ErrInvalidCredentials,
			accessTokenExpected:  false,
			refreshTokenExpected: false,
		},
		{
			name:      "error: inactive user",
			email:     "user@critiquefi.com",
			password:  "password",
			userAgent: "test",
			remember:  true,
			storeSetup: func(s *mockStore) {
				s.On(GetUserByEmail, models.User{
					ID:           1,
					Email:        "user@critiquefi.com",
					DisplayName:  "test.user",
					Name:         "Test User",
					IsAdmin:      false,
					PasswordHash: "wrongpassword",
					IsActive:     false,
				}, nil)
			},
			expectedErr:          business.ErrInvalidCredentials,
			accessTokenExpected:  false,
			refreshTokenExpected: false,
		},
	}

	for _, tc := range tc {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := newMockStore(t)
			tc.storeSetup(s)

			svc := authbus.New(authbus.BusConfig{
				Store:           s,
				AccessTokenKey:  "key",
				AccessTokenTTL:  time.Hour,
				RefreshTokenTTL: time.Hour,
			})

			accessToken, refreshToken, err := svc.Login(context.Background(), tc.email, tc.password, tc.userAgent, tc.remember)
			testutil.CheckErr(err, tc.expectedErr, t)

			if tc.accessTokenExpected && accessToken == "" {
				t.Error("expected access token, got nothing")
			}
			if !tc.accessTokenExpected && accessToken != "" {
				t.Errorf("expected no access token, got %v", accessToken)
			}
			if tc.refreshTokenExpected && refreshToken == "" {
				t.Error("expected refresh token, got nothing")
			}
			if !tc.refreshTokenExpected && refreshToken != "" {
				t.Errorf("expected no refresh token, got %v", refreshToken)
			}
		})
	}
}

func TestBus_Logout(t *testing.T) {
	t.Parallel()
	tc := []struct {
		name         string
		refreshToken string
		storeSetup   func(s *mockStore)
		expectedErr  error
	}{
		{
			name:         "success: logged out",
			refreshToken: "test-refresh-token",
			storeSetup: func(s *mockStore) {
				s.On(DeleteRefreshToken, nil)
			},
			expectedErr: nil,
		},
		{
			name:         "success: no error if token not found",
			refreshToken: "test-refresh-token",
			storeSetup: func(s *mockStore) {
				s.On(DeleteRefreshToken, store.ErrNotFound)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range tc {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := newMockStore(t)
			tc.storeSetup(s)

			svc := authbus.New(authbus.BusConfig{
				Store:           s,
				AccessTokenKey:  "key",
				AccessTokenTTL:  time.Hour,
				RefreshTokenTTL: time.Hour,
			})

			err := svc.Logout(context.Background(), tc.refreshToken)
			testutil.CheckErr(err, tc.expectedErr, t)
		})
	}
}

func TestBus_Refresh(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name         string
		refreshToken string
		storeSetup   func(s *mockStore)
		expectedErr  error
	}{
		{
			name:         "success: refreshed",
			refreshToken: "refresh_token",
			storeSetup: func(s *mockStore) {
				s.On(GetRefreshToken, models.RefreshToken{
					TokenHash: "refresh_token",
					UserID:    1,
					ExpiresAt: time.Now().Add(time.Hour),
					CreatedAt: time.Now(),
				}, nil).
					On(DeleteRefreshToken, nil).
					On(GetUserByID, models.User{
						ID:           1,
						Email:        "user@critiquefi.com",
						DisplayName:  "test.user",
						Name:         "Test User",
						IsAdmin:      false,
						PasswordHash: "password",
						IsActive:     true,
					}, nil).
					On(CreateRefreshToken, nil)
			},
			expectedErr: nil,
		},
		{
			name:         "error: invalid refresh token - not found",
			refreshToken: "refresh_token",
			storeSetup: func(s *mockStore) {
				s.On(GetRefreshToken, models.RefreshToken{}, store.ErrNotFound)
			},
			expectedErr: business.ErrInvalidToken,
		},
		{
			name:         "error: invalid refresh token - expired",
			refreshToken: "refresh_token",
			storeSetup: func(s *mockStore) {
				s.On(GetRefreshToken, models.RefreshToken{
					TokenHash: "refresh_token",
					UserID:    1,
					ExpiresAt: time.Now().Add(-time.Hour),
					CreatedAt: time.Now().Add(-time.Hour),
				}, nil).
					On(DeleteRefreshToken, nil)
			},
			expectedErr: business.ErrInvalidToken,
		},
		{
			name:         "error: invalid refresh token - user not found",
			refreshToken: "refresh_token",
			storeSetup: func(s *mockStore) {
				s.On(GetRefreshToken, models.RefreshToken{
					TokenHash: "refresh_token",
					UserID:    1,
					ExpiresAt: time.Now().Add(time.Hour),
					CreatedAt: time.Now(),
				}, nil).
					On(DeleteRefreshToken, nil).
					On(GetUserByID, models.User{}, store.ErrNotFound)
			},
			expectedErr: business.ErrInvalidToken,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := newMockStore(t)
			tc.storeSetup(s)

			svc := authbus.New(authbus.BusConfig{
				Store:           s,
				AccessTokenKey:  "key",
				AccessTokenTTL:  time.Hour,
				RefreshTokenTTL: time.Hour,
			})

			accessToken, newRefreshToken, err := svc.Refresh(context.Background(), tc.refreshToken)
			testutil.CheckErr(err, tc.expectedErr, t)

			if tc.expectedErr != nil {
				return
			}

			if len(accessToken) == 0 {
				t.Error("expected access token, got nothing")
			}
			if len(newRefreshToken) == 0 {
				t.Error("expected refresh token, got nothing")
			}
			if tc.refreshToken == newRefreshToken {
				t.Error("expected rotated refresh token, got the same one")
			}
		})
	}
}

func TestBus_ValidateAccessToken(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name          string
		generateToken func(s *authbus.Bus) string
		expectedErr   error
		checkClaims   func(claims models.Claims)
	}{
		{
			name: "success - valid user token",
			generateToken: func(s *authbus.Bus) string {
				token, err := s.GenerateAccessToken(models.User{
					ID:           1,
					Email:        "user@critiquefi.com",
					DisplayName:  "test.user",
					Name:         "Test User",
					IsAdmin:      false,
					PasswordHash: "password",
					IsActive:     true,
				})
				if err != nil {
					t.Fatal(err)
				}
				return token
			},
			expectedErr: nil,
			checkClaims: func(claims models.Claims) {
				if claims.UserID != int64(1) {
					t.Errorf("expected user id 1, got %v", claims.UserID)
				}
				if claims.IsAdmin {
					t.Errorf("expected is admin to be false, got %v", claims.IsAdmin)
				}
				if claims.Email != "user@critiquefi.com" {
					t.Errorf("expected mail user@critiquefi.com, got %v", claims.Email)
				}
			},
		},
		{
			name: "success - valid admin token",
			generateToken: func(s *authbus.Bus) string {
				token, err := s.GenerateAccessToken(models.User{
					ID:           1,
					Email:        "user@critiquefi.com",
					DisplayName:  "test.user",
					Name:         "Test User",
					IsAdmin:      true,
					PasswordHash: "password",
					IsActive:     true,
				})
				if err != nil {
					t.Fatal(err)
				}
				return token
			},
			expectedErr: nil,
			checkClaims: func(claims models.Claims) {
				if claims.UserID != int64(1) {
					t.Errorf("expected user id 1, got %v", claims.UserID)
				}
				if !claims.IsAdmin {
					t.Errorf("expected is admin to be true, got %v", claims.IsAdmin)
				}
				if claims.Email != "user@critiquefi.com" {
					t.Errorf("expected mail user@critiquefi.com, got %v", claims.Email)
				}
			},
		},
		{
			name: "error - invalid token - empty",
			generateToken: func(s *authbus.Bus) string {
				return ""
			},
			expectedErr: business.ErrInvalidToken,
		},
		{
			name: "error - invalid token - malformed",
			generateToken: func(s *authbus.Bus) string {
				return "WHAT THE!! WHAT IS THIS???"
			},
			expectedErr: business.ErrInvalidToken,
		},
		{
			name: "error - invalid token - wrong key",
			generateToken: func(s *authbus.Bus) string {
				wrongService := authbus.New(authbus.BusConfig{
					AccessTokenKey: "wrongkey",
					AccessTokenTTL: time.Hour,
				})
				token, err := wrongService.GenerateAccessToken(models.User{
					ID:           1,
					Email:        "user@critiquefi.com",
					DisplayName:  "test.user",
					Name:         "Test User",
					IsAdmin:      true,
					PasswordHash: "password",
					IsActive:     true,
				})
				if err != nil {
					t.Fatal(err)
				}
				return token
			},
			expectedErr: business.ErrInvalidToken,
		},
		{
			name: "error - invalid token - expired",
			generateToken: func(s *authbus.Bus) string {
				expiredService := authbus.New(authbus.BusConfig{
					AccessTokenKey: "key",
					AccessTokenTTL: -1 * time.Hour,
				})
				token, err := expiredService.GenerateAccessToken(models.User{
					ID:           1,
					Email:        "user@critiquefi.com",
					DisplayName:  "test.user",
					Name:         "Test User",
					IsAdmin:      false,
					PasswordHash: "password",
					IsActive:     true,
				})
				if err != nil {
					t.Fatal(err)
				}
				return token
			},
			expectedErr: business.ErrInvalidToken,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			svc := authbus.New(authbus.BusConfig{
				AccessTokenKey: "key",
				AccessTokenTTL: time.Hour,
			})

			token := tc.generateToken(svc)
			claims, err := svc.ValidateAccessToken(token)
			testutil.CheckErr(err, tc.expectedErr, t)
			if tc.expectedErr != nil {
				return
			}
			tc.checkClaims(claims)
		})
	}
}
