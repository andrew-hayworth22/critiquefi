package syshttp_test

import (
	"context"
	"testing"

	"github.com/andrew-hayworth22/critiquefi-service/internal/testutil"
)

var (
	Ping testutil.Method = "Ping"
)

// mockBus is a mock sysbus bus for testing
type mockBus struct {
	testutil.Mock
}

func newMockBus(t *testing.T) *mockBus {
	return &mockBus{Mock: testutil.NewMock(t)}
}

func (b *mockBus) Ping(ctx context.Context) error {
	call := b.Next(Ping)
	return testutil.ConvertError(call.Returns[0])
}
