package syshttp_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andrew-hayworth22/critiquefi-service/internal/business"
	"github.com/andrew-hayworth22/critiquefi-service/internal/http/syshttp"
	"github.com/andrew-hayworth22/critiquefi-service/internal/testutil"
)

func TestHandler_Liveness(t *testing.T) {
	t.Parallel()

	bus := newMockBus(t)

	handler := syshttp.New(bus)
	req := testutil.NewTestRequest(t, http.MethodGet, "/liveness", nil)
	w := httptest.NewRecorder()

	handler.Liveness(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var resp syshttp.SystemCheckResponse
	testutil.DecodeResponse(t, w, &resp)
	if resp.Status != "ok" {
		t.Errorf("expected status to be ok, got %s", resp.Status)
	}
}

func TestHandler_Readiness(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name           string
		busSetup       func(s *mockBus)
		expectedStatus int
		checkResponse  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "success: ready",
			busSetup: func(s *mockBus) {
				s.On(Ping, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp syshttp.SystemCheckResponse
				testutil.DecodeResponse(t, w, &resp)
				if resp.Status != "ok" {
					t.Errorf("expected status to be ok, got %s", resp.Status)
				}
			},
		},
		{
			name: "error: ping failed",
			busSetup: func(s *mockBus) {
				s.On(Ping, business.ErrInternal)
			},
			expectedStatus: http.StatusServiceUnavailable,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			bus := newMockBus(t)
			tc.busSetup(bus)

			handler := syshttp.New(bus)
			req := testutil.NewTestRequest(t, http.MethodGet, "/readiness", nil)
			w := httptest.NewRecorder()

			handler.Readiness(w, req)
			if w.Code != tc.expectedStatus {
				t.Errorf("expected status code %d, got %d", tc.expectedStatus, w.Code)
			}

			if tc.checkResponse != nil {
				tc.checkResponse(t, w)
			}
		})
	}
}
