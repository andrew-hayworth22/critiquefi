package authhttp_test

import (
	"context"
	"testing"

	"github.com/andrew-hayworth22/critiquefi/services/api/internal/models"
	"github.com/andrew-hayworth22/critiquefi/services/api/internal/testutil"
)

var (
	Register       testutil.Method = "Register"
	Login          testutil.Method = "Login"
	Logout         testutil.Method = "Logout"
	Refresh        testutil.Method = "Refresh"
	ForgotPassword testutil.Method = "ForgotPassword"
	ResetPassword  testutil.Method = "ResetPassword"
)

// mockBus is a mock authbus for testing
type mockBus struct {
	testutil.Mock
}

func newMockBus(t *testing.T) *mockBus {
	return &mockBus{Mock: testutil.NewMock(t)}
}

func (b *mockBus) Register(ctx context.Context, user models.NewUserRequest, userAgent string, remember bool) (string, string, error) {
	call := b.Next(Register)
	return call.Returns[0].(string), call.Returns[1].(string), testutil.ConvertError(call.Returns[2])
}

func (b *mockBus) Login(ctx context.Context, email, password string, userAgent string, remember bool) (string, string, error) {
	call := b.Next(Login)
	return call.Returns[0].(string), call.Returns[1].(string), testutil.ConvertError(call.Returns[2])
}

func (b *mockBus) Logout(ctx context.Context, refreshToken string) error {
	call := b.Next(Logout)
	return testutil.ConvertError(call.Returns[0])
}

func (b *mockBus) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	call := b.Next(Refresh)
	return call.Returns[0].(string), call.Returns[1].(string), testutil.ConvertError(call.Returns[2])
}

func (b *mockBus) ForgotPassword(ctx context.Context, email string) error {
	call := b.Next(ForgotPassword)
	return testutil.ConvertError(call.Returns[0])
}

func (b *mockBus) ResetPassword(ctx context.Context, token, password string) error {
	call := b.Next(ResetPassword)
	return testutil.ConvertError(call.Returns[0])
}
