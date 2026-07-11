package sysbus_test

import (
	"context"
	"testing"

	"github.com/andrew-hayworth22/critiquefi-service/internal/testutil"
)

var (
	Ping testutil.Method = "Ping"
)

// mockStore is a mock sysbus store for testing
type mockStore struct {
	testutil.Mock
}

func newMockStore(t *testing.T) *mockStore {
	t.Helper()

	return &mockStore{Mock: testutil.NewMock(t)}
}

func (s *mockStore) Ping(ctx context.Context) error {
	call := s.Next(Ping)
	return testutil.ConvertError(call.Returns[0])
}
