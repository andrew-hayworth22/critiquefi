package authbus_test

import (
	"context"
	"testing"

	"github.com/andrew-hayworth22/critiquefi/services/api/internal/mail"
	"github.com/andrew-hayworth22/critiquefi/services/api/internal/models"
	"github.com/andrew-hayworth22/critiquefi/services/api/internal/testutil"
)

var (
	CreateUser           testutil.Method = "CreateUser"
	GetUserByID          testutil.Method = "GetUserByID"
	GetUserByEmail       testutil.Method = "GetUserByEmail"
	CheckTakenUserFields testutil.Method = "CheckTakenUserFields"
	SetUserLastLogin     testutil.Method = "SetUserLastLogin"
	SetUserPassword      testutil.Method = "SetUserPassword"

	CreateRefreshToken testutil.Method = "CreateRefreshToken"
	GetRefreshToken    testutil.Method = "GetRefreshToken"
	DeleteRefreshToken testutil.Method = "DeleteRefreshToken"

	CreatePasswordResetToken          testutil.Method = "CreatePasswordResetToken"
	GetPasswordResetToken             testutil.Method = "GetPasswordResetToken"
	DeletePasswordResetTokensByUserID testutil.Method = "DeletePasswordResetTokensByUserID"

	SendWelcome       testutil.Method = "SendWelcome"
	SendPasswordReset testutil.Method = "SendPasswordReset"
)

// mockStore is a mock authbus store for testing
type mockStore struct {
	testutil.Mock
}

func newMockStore(t *testing.T) *mockStore {
	t.Helper()

	return &mockStore{
		Mock: testutil.NewMock(t),
	}
}

func (s *mockStore) CreateUser(ctx context.Context, user models.NewUser) (int64, error) {
	call := s.Next(CreateUser)
	return call.Returns[0].(int64), testutil.ConvertError(call.Returns[1])
}

func (s *mockStore) GetUserByID(ctx context.Context, id int64) (models.User, error) {
	call := s.Next(GetUserByID)
	return call.Returns[0].(models.User), testutil.ConvertError(call.Returns[1])
}

func (s *mockStore) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	call := s.Next(GetUserByEmail)
	return call.Returns[0].(models.User), testutil.ConvertError(call.Returns[1])
}

func (s *mockStore) CheckTakenUserFields(ctx context.Context, request models.NewUserRequest) (models.UserFieldsTaken, error) {
	call := s.Next(CheckTakenUserFields)
	return call.Returns[0].(models.UserFieldsTaken), testutil.ConvertError(call.Returns[1])
}

func (s *mockStore) SetUserLastLogin(ctx context.Context, id int64) error {
	call := s.Next(SetUserLastLogin)
	return testutil.ConvertError(call.Returns[0])
}

func (s *mockStore) SetUserPassword(ctx context.Context, id int64, password string) error {
	call := s.Next(SetUserPassword)
	return testutil.ConvertError(call.Returns[0])
}

func (s *mockStore) CreateRefreshToken(ctx context.Context, refreshToken models.RefreshToken) error {
	call := s.Next(CreateRefreshToken)
	return testutil.ConvertError(call.Returns[0])
}

func (s *mockStore) GetRefreshToken(ctx context.Context, tokenHash string) (models.RefreshToken, error) {
	call := s.Next(GetRefreshToken)
	return call.Returns[0].(models.RefreshToken), testutil.ConvertError(call.Returns[1])
}

func (s *mockStore) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	call := s.Next(DeleteRefreshToken)
	return testutil.ConvertError(call.Returns[0])
}

func (s *mockStore) CreatePasswordResetToken(ctx context.Context, token models.PasswordResetToken) error {
	call := s.Next(CreatePasswordResetToken)
	return testutil.ConvertError(call.Returns[0])
}

func (s *mockStore) GetPasswordResetToken(ctx context.Context, tokenHash string) (models.PasswordResetToken, error) {
	call := s.Next(GetPasswordResetToken)
	return call.Returns[0].(models.PasswordResetToken), testutil.ConvertError(call.Returns[1])
}

func (s *mockStore) DeletePasswordResetTokensByUserID(ctx context.Context, userID int64) error {
	call := s.Next(DeletePasswordResetTokensByUserID)
	return testutil.ConvertError(call.Returns[0])
}

// mockMailer is a mock auth mailer for testing
type mockMailer struct {
	testutil.Mock
}

func newMockMailer(t *testing.T) *mockMailer {
	t.Helper()

	return &mockMailer{
		Mock: testutil.NewMock(t),
	}
}

func (m *mockMailer) SendWelcome(ctx context.Context, recipient mail.Recipient) error {
	call := m.Next(SendWelcome)
	return testutil.ConvertError(call.Returns[0])
}

func (m *mockMailer) SendPasswordReset(ctx context.Context, recipient mail.Recipient, token string) error {
	call := m.Next(SendPasswordReset)
	return testutil.ConvertError(call.Returns[0])
}
