package sysbus_test

import (
	"context"
	"testing"

	"github.com/andrew-hayworth22/critiquefi-service/internal/business"
	"github.com/andrew-hayworth22/critiquefi-service/internal/business/sysbus"
	"github.com/andrew-hayworth22/critiquefi-service/internal/store"
	"github.com/andrew-hayworth22/critiquefi-service/internal/testutil"
)

func TestBus_Ping(t *testing.T) {
	t.Run("ping", func(t *testing.T) {
		t.Parallel()
		tc := []struct {
			name        string
			storeSetup  func(s *mockStore)
			expectedErr error
		}{
			{
				name: "success: pong",
				storeSetup: func(s *mockStore) {
					s.On(Ping, nil)
				},
				expectedErr: nil,
			},
			{
				name: "error: ping failed",
				storeSetup: func(s *mockStore) {
					s.On(Ping, store.ErrInternal)
				},
				expectedErr: business.ErrInternal,
			},
		}

		for _, tc := range tc {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				s := newMockStore(t)
				tc.storeSetup(s)

				bus := sysbus.New(s)

				err := bus.Ping(context.Background())
				testutil.CheckErr(err, tc.expectedErr, t)
			})
		}
	})
}
