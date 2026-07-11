package authpg_test

import (
	"context"
	"testing"
	"time"

	"github.com/andrew-hayworth22/critiquefi-service/internal/models"
	"github.com/andrew-hayworth22/critiquefi-service/internal/store"
	"github.com/andrew-hayworth22/critiquefi-service/internal/store/postgres/authpg"
	"github.com/andrew-hayworth22/critiquefi-service/internal/testutil"
)

func TestAuthStore_CreateUser(t *testing.T) {
	t.Run("create user", func(t *testing.T) {
		testutil.PrepareDB(t, testDB)
		s := authpg.New(testDB)

		cases := []struct {
			name        string
			user        models.NewUser
			expectedErr error
		}{
			{
				name: "success: creates user",
				user: models.NewUser{
					Email:        "test@example.com",
					DisplayName:  "test.example",
					Name:         "Test User",
					PasswordHash: "hashedpassword",
				},
				expectedErr: nil,
			},
			{
				name: "error: duplicate email",
				user: models.NewUser{
					Email:        "test@example.com",
					DisplayName:  "test.example2",
					Name:         "Test User",
					PasswordHash: "hashedpassword",
				},
				expectedErr: store.ErrDuplicate,
			},
			{
				name: "error: duplicate display name",
				user: models.NewUser{
					Email:        "test2@example.com",
					DisplayName:  "test.example",
					Name:         "Test User",
					PasswordHash: "hashedpassword",
				},
				expectedErr: store.ErrDuplicate,
			},
			{
				name: "error: duplicate email and display name",
				user: models.NewUser{
					Email:        "test@example.com",
					DisplayName:  "test.example",
					Name:         "Test User",
					PasswordHash: "hashedpassword",
				},
				expectedErr: store.ErrDuplicate,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				id, err := s.CreateUser(context.Background(), tc.user)
				testutil.CheckErr(err, tc.expectedErr, t)

				if tc.expectedErr == nil && id == 0 {
					t.Fatal("expected non-zero user ID")
				}
			})
		}
	})

}

func TestAuthStore_GetUserById(t *testing.T) {
	t.Run("get user by ID", func(t *testing.T) {
		testutil.PrepareDB(t, testDB)
		s := authpg.New(testDB)

		cases := []struct {
			name         string
			userID       int64
			expectedUser models.User
			expectedErr  error
		}{
			{
				name:   "success: gets user",
				userID: 1,
				expectedUser: models.User{
					ID:           1,
					Email:        "user@critiquefi.com",
					DisplayName:  "test.user",
					Name:         "Test User",
					IsAdmin:      false,
					PasswordHash: "password",
					IsActive:     true,
				},
				expectedErr: nil,
			},
			{
				name:   "success: gets admin",
				userID: 2,
				expectedUser: models.User{
					ID:           2,
					Email:        "admin@critiquefi.com",
					DisplayName:  "admin.user",
					Name:         "Test Admin",
					IsAdmin:      true,
					PasswordHash: "password",
					IsActive:     true,
				},
				expectedErr: nil,
			},
			{
				name:   "success: gets inactive user",
				userID: 3,
				expectedUser: models.User{
					ID:           3,
					Email:        "deactivated@critiquefi.com",
					DisplayName:  "deactivated.user",
					Name:         "Deactivated User",
					IsAdmin:      false,
					PasswordHash: "password",
					IsActive:     false,
				},
				expectedErr: nil,
			},
			{
				name:         "error: user not found",
				userID:       4,
				expectedUser: models.User{},
				expectedErr:  store.ErrNotFound,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				user, err := s.GetUserByID(context.Background(), tc.userID)
				testutil.CheckErr(err, tc.expectedErr, t)

				if user != tc.expectedUser {
					t.Errorf("expected user: %v, got: %v", tc.expectedUser, user)
				}
			})
		}
	})
}

func TestAuthStore_GetUserByEmail(t *testing.T) {
	t.Run("get user by email", func(t *testing.T) {
		testutil.PrepareDB(t, testDB)
		s := authpg.New(testDB)

		cases := []struct {
			name         string
			email        string
			expectedUser models.User
			expectedErr  error
		}{
			{
				name:  "success: gets user",
				email: "user@critiquefi.com",
				expectedUser: models.User{
					ID:           1,
					Email:        "user@critiquefi.com",
					DisplayName:  "test.user",
					Name:         "Test User",
					IsAdmin:      false,
					PasswordHash: "password",
					IsActive:     true,
				},
				expectedErr: nil,
			},
			{
				name:  "success: gets admin",
				email: "admin@critiquefi.com",
				expectedUser: models.User{
					ID:           2,
					Email:        "admin@critiquefi.com",
					DisplayName:  "admin.user",
					Name:         "Test Admin",
					IsAdmin:      true,
					PasswordHash: "password",
					IsActive:     true,
				},
				expectedErr: nil,
			},
			{
				name:  "success: gets inactive user",
				email: "deactivated@critiquefi.com",
				expectedUser: models.User{
					ID:           3,
					Email:        "deactivated@critiquefi.com",
					DisplayName:  "deactivated.user",
					Name:         "Deactivated User",
					IsAdmin:      false,
					PasswordHash: "password",
					IsActive:     false,
				},
				expectedErr: nil,
			},
			{
				name:         "error: user not found",
				email:        "nonexistent@critiquefi.com",
				expectedUser: models.User{},
				expectedErr:  store.ErrNotFound,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				user, err := s.GetUserByEmail(context.Background(), tc.email)
				testutil.CheckErr(err, tc.expectedErr, t)

				if user != tc.expectedUser {
					t.Errorf("expected user: %v, got: %v", tc.expectedUser, user)
				}
			})
		}
	})
}

func TestAuthStore_CheckTakenFields(t *testing.T) {
	t.Run("check taken fields", func(t *testing.T) {
		testutil.PrepareDB(t, testDB)
		s := authpg.New(testDB)

		cases := []struct {
			name                    string
			newUser                 models.NewUserRequest
			expectedUserFieldsTaken models.UserFieldsTaken
			expectedErr             error
		}{
			{
				name: "success: no fields taken",
				newUser: models.NewUserRequest{
					Email:           "new@critiquefi.com",
					DisplayName:     "new.user",
					Name:            "New User",
					Password:        "password",
					ConfirmPassword: "password",
				},
				expectedUserFieldsTaken: models.UserFieldsTaken{
					EmailTaken:       false,
					DisplayNameTaken: false,
				},
				expectedErr: nil,
			},
			{
				name: "success: email taken",
				newUser: models.NewUserRequest{
					Email:           "user@critiquefi.com",
					DisplayName:     "test.user2",
					Name:            "Test User",
					Password:        "password",
					ConfirmPassword: "password",
				},
				expectedUserFieldsTaken: models.UserFieldsTaken{
					EmailTaken:       true,
					DisplayNameTaken: false,
				},
				expectedErr: nil,
			},
			{
				name: "success: display name taken",
				newUser: models.NewUserRequest{
					Email:           "user2@critiquefi.com",
					DisplayName:     "test.user",
					Name:            "Test User",
					Password:        "password",
					ConfirmPassword: "password",
				},
				expectedUserFieldsTaken: models.UserFieldsTaken{
					EmailTaken:       false,
					DisplayNameTaken: true,
				},
				expectedErr: nil,
			},
			{
				name: "success: display name and email taken",
				newUser: models.NewUserRequest{
					Email:           "user@critiquefi.com",
					DisplayName:     "test.user",
					Name:            "Test User",
					Password:        "password",
					ConfirmPassword: "password",
				},
				expectedUserFieldsTaken: models.UserFieldsTaken{
					EmailTaken:       true,
					DisplayNameTaken: true,
				},
				expectedErr: nil,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				userFieldsTaken, err := s.CheckTakenUserFields(context.Background(), tc.newUser)
				testutil.CheckErr(err, tc.expectedErr, t)

				if userFieldsTaken != tc.expectedUserFieldsTaken {
					t.Errorf("expected user fields taken: %v, got: %v", tc.expectedUserFieldsTaken, userFieldsTaken)
				}
			})
		}
	})
}

func TestAuthStore_CreateRefreshToken(t *testing.T) {
	t.Run("create refresh token", func(t *testing.T) {
		testutil.PrepareDB(t, testDB)
		s := authpg.New(testDB)

		cases := []struct {
			name         string
			refreshToken models.RefreshToken
			expectedErr  error
		}{
			{
				name: "success: creates refresh token",
				refreshToken: models.RefreshToken{
					TokenHash: "tokenhash",
					UserID:    1,
					UserAgent: "useragent",
					ExpiresAt: time.Now(),
					CreatedAt: time.Now(),
				},
				expectedErr: nil,
			},
			{
				name: "error: duplicate hash",
				refreshToken: models.RefreshToken{
					TokenHash: "tokenhash",
					UserID:    1,
					UserAgent: "useragent",
					ExpiresAt: time.Now(),
					CreatedAt: time.Now(),
				},
				expectedErr: store.ErrDuplicate,
			},
			{
				name: "error: foreign key violation",
				refreshToken: models.RefreshToken{
					TokenHash: "tokenhash2",
					UserID:    0,
					UserAgent: "useragent",
					ExpiresAt: time.Now(),
					CreatedAt: time.Now(),
				},
				expectedErr: store.ErrForeignKeyViolation,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				err := s.CreateRefreshToken(context.Background(), tc.refreshToken)
				testutil.CheckErr(err, tc.expectedErr, t)
			})
		}
	})
}

func TestAuthStore_GetRefreshToken(t *testing.T) {
	t.Run("get refresh token", func(t *testing.T) {
		testutil.PrepareDB(t, testDB)
		s := authpg.New(testDB)

		cases := []struct {
			name          string
			tokenHash     string
			expectedToken models.RefreshToken
			expectedErr   error
		}{
			{
				name:      "success: gets refresh token",
				tokenHash: "test-token-hash",
				expectedToken: models.RefreshToken{
					TokenHash: "test-token-hash",
					UserID:    1,
					UserAgent: "test-user-agent",
					ExpiresAt: time.Time{},
					CreatedAt: time.Time{},
				},
				expectedErr: nil,
			},
			{
				name:          "error: token not found",
				tokenHash:     "test-token-hash2",
				expectedToken: models.RefreshToken{},
				expectedErr:   store.ErrNotFound,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				token, err := s.GetRefreshToken(context.Background(), tc.tokenHash)
				testutil.CheckErr(err, tc.expectedErr, t)

				if token.TokenHash != tc.expectedToken.TokenHash {
					t.Errorf("expected token hash: %v, got: %v", tc.expectedToken.TokenHash, token.TokenHash)
				}
			})
		}
	})
}

func TestAuthStore_DeleteRefreshToken(t *testing.T) {
	t.Run("delete refresh token", func(t *testing.T) {
		testutil.PrepareDB(t, testDB)
		s := authpg.New(testDB)

		cases := []struct {
			name        string
			tokenHash   string
			expectedErr error
		}{
			{
				name:        "success: deletes refresh token",
				tokenHash:   "test-token-hash",
				expectedErr: nil,
			},
			{
				name:        "success: no error if token not found",
				tokenHash:   "test-token-hash2",
				expectedErr: nil,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				err := s.DeleteRefreshToken(context.Background(), tc.tokenHash)
				testutil.CheckErr(err, tc.expectedErr, t)

				rows, err := testDB.QueryContext(context.Background(), "SELECT * FROM refresh_tokens WHERE token_hash = $1", tc.tokenHash)
				if err != nil || rows.Next() {
					t.Errorf("expected no rows, got: %v", err)
				}
			})
		}
	})
}

func TestAuthStore_SetUserLastLogin(t *testing.T) {
	t.Run("set user last login", func(t *testing.T) {
		testutil.PrepareDB(t, testDB)
		s := authpg.New(testDB)

		cases := []struct {
			name   string
			userID int64
		}{
			{
				name:   "success: sets last login",
				userID: 1,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				getLastLogin := func() time.Time {
					q := "SELECT last_login FROM users WHERE id = $1"

					row, err := testDB.QueryContext(context.Background(), q, tc.userID)
					if err != nil {
						t.Errorf("expected no error, got: %v", err)
					}
					defer row.Close()
					if !row.Next() {
						t.Error("expected row, got nothing")
					}
					var lastLogin time.Time
					if err := row.Scan(&lastLogin); err != nil {
						t.Errorf("expected no error, got: %v", err)
					}
					return lastLogin
				}

				initialLastLogin := getLastLogin()

				err := s.SetUserLastLogin(context.Background(), tc.userID)
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}

				lastLogin := getLastLogin()

				if lastLogin.Equal(initialLastLogin) {
					t.Errorf("expected last login to be updated, got: %v", lastLogin)
				}
			})
		}
	})
}
