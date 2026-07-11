package authhttp_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andrew-hayworth22/critiquefi-service/internal/business"
	"github.com/andrew-hayworth22/critiquefi-service/internal/http/authhttp"
	"github.com/andrew-hayworth22/critiquefi-service/internal/models"
	"github.com/andrew-hayworth22/critiquefi-service/internal/testutil"
)

func TestHandler_Register(t *testing.T) {
	t.Parallel()

	cookieName := "rt"

	cases := []struct {
		name           string
		body           any
		busSetup       func(s *mockBus)
		expectedStatus int
		checkResponse  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "success: registered user (don't remember)",
			body: map[string]any{
				"email":            "user@critiquefi.com",
				"display_name":     "test.user",
				"name":             "Test User",
				"password":         "password",
				"confirm_password": "password",
				"remember":         false,
			},
			busSetup: func(s *mockBus) {
				s.On(Register, "access-token", "", nil)
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response authhttp.AuthenticationResponse
				testutil.DecodeResponse(t, w, &response)

				if response.AccessToken != "access-token" {
					t.Errorf("expected access token to be 'access-token', got %s", response.AccessToken)
				}

				cookies := w.Result().Cookies()
				if len(cookies) > 0 {
					t.Errorf("expected no cookies, got %d", len(cookies))
				}
			},
		},
		{
			name: "success: registered user (remember)",
			body: map[string]any{
				"email":            "user@critiquefi.com",
				"display_name":     "test.user",
				"name":             "Test User",
				"password":         "password",
				"confirm_password": "password",
				"remember":         true,
			},
			busSetup: func(s *mockBus) {
				s.On(Register, "access-token", "refresh-token", nil)
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response authhttp.AuthenticationResponse
				testutil.DecodeResponse(t, w, &response)

				if response.AccessToken != "access-token" {
					t.Errorf("expected access token to be 'access-token', got %s", response.AccessToken)
				}

				cookies := w.Result().Cookies()
				if len(cookies) != 1 {
					t.Errorf("expected 1 cookie, got %d", len(cookies))
				}
				if cookies[0].Name != cookieName {
					t.Errorf("expected cookie name to be %s, got %s", cookieName, cookies[0].Name)
				}
				if cookies[0].Value != "refresh-token" {
					t.Errorf("expected cookie value to be 'refresh-token', got %s", cookies[0].Value)
				}
			},
		},
		{
			name: "error: validation error",
			body: map[string]any{
				"email":            "user@critiquefi.com",
				"display_name":     "test.user",
				"name":             "Test User",
				"password":         "password",
				"confirm_password": "password",
				"remember":         false,
			},
			busSetup: func(s *mockBus) {
				s.On(Register, "", "", business.ErrDuplicate)
			},
			expectedStatus: http.StatusConflict,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				testutil.DecodeResponse(t, w, &response)

				if len(response) != 1 {
					t.Errorf("expected 1 error, got %d", len(response))
				}

				cookies := w.Result().Cookies()
				if len(cookies) != 0 {
					t.Errorf("expected 0 cookies, got %d", len(cookies))
				}
			},
		},
		{
			name: "error: duplicate error",
			body: map[string]any{
				"email":            "user@critiquefi.com",
				"display_name":     "test.user",
				"name":             "Test User",
				"password":         "password",
				"confirm_password": "password",
				"remember":         false,
			},
			busSetup: func(s *mockBus) {
				s.On(Register, "", "", models.ValidationErrors{
					"email":        "invalid email address",
					"display_name": "display name must be between 3 and 20 characters long",
				})
			},
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				testutil.DecodeResponse(t, w, &response)

				if len(response) != 2 {
					t.Errorf("expected 2 errors, got %d", len(response))
				}

				cookies := w.Result().Cookies()
				if len(cookies) != 0 {
					t.Errorf("expected 0 cookies, got %d", len(cookies))
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			bus := newMockBus(t)
			tc.busSetup(bus)

			handler := authhttp.New(bus, cookieName, "critiquefi.com")
			req := testutil.NewTestRequest(t, http.MethodPost, "/register", tc.body)
			w := httptest.NewRecorder()

			handler.Register(w, req)
			if w.Code != tc.expectedStatus {
				t.Errorf("expected status code %d, got %d", tc.expectedStatus, w.Code)
			}

			if tc.checkResponse != nil {
				tc.checkResponse(t, w)
			}
		})
	}
}

func TestHandler_Login(t *testing.T) {
	t.Parallel()

	cookieName := "rt"

	cases := []struct {
		name           string
		body           any
		busSetup       func(s *mockBus)
		expectedStatus int
		checkResponse  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "success: logged in user (don't remember)",
			body: map[string]any{
				"email":    "user@critiquefi.com",
				"password": "password",
				"remember": false,
			},
			busSetup: func(s *mockBus) {
				s.On(Login, "access-token", "", nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response authhttp.AuthenticationResponse
				testutil.DecodeResponse(t, w, &response)

				if response.AccessToken != "access-token" {
					t.Errorf("expected access token to be 'access-token', got %s", response.AccessToken)
				}

				cookies := w.Result().Cookies()
				if len(cookies) > 0 {
					t.Errorf("expected no cookies, got %d", len(cookies))
				}
			},
		},
		{
			name: "success: logged in user (remember)",
			body: map[string]any{
				"email":    "user@critiquefi.com",
				"password": "password",
				"remember": true,
			},
			busSetup: func(s *mockBus) {
				s.On(Login, "access-token", "refresh-token", nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response authhttp.AuthenticationResponse
				testutil.DecodeResponse(t, w, &response)

				if response.AccessToken != "access-token" {
					t.Errorf("expected access token to be 'access-token', got %s", response.AccessToken)
				}

				cookies := w.Result().Cookies()
				if len(cookies) != 1 {
					t.Errorf("expected 1 cookie, got %d", len(cookies))
				}
				if cookies[0].Name != cookieName {
					t.Errorf("expected cookie name to be %s, got %s", cookieName, cookies[0].Name)
				}
				if cookies[0].Value != "refresh-token" {
					t.Errorf("expected cookie value to be 'refresh-token', got %s", cookies[0].Value)
				}
			},
		},
		{
			name: "error: invalid credentials",
			body: map[string]any{
				"email":    "user@critiquefi.com",
				"password": "password",
				"remember": true,
			},
			busSetup: func(s *mockBus) {
				s.On(Login, "", "", business.ErrInvalidCredentials)
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				testutil.DecodeResponse(t, w, &response)

				if len(response) != 1 {
					t.Errorf("expected 1 error, got %v", len(response))
				}

				cookies := w.Result().Cookies()
				if len(cookies) != 0 {
					t.Errorf("expected 0 cookies, got %d", len(cookies))
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			bus := newMockBus(t)
			tc.busSetup(bus)

			handler := authhttp.New(bus, cookieName, "critiquefi.com")
			req := testutil.NewTestRequest(t, http.MethodPost, "/register", tc.body)
			w := httptest.NewRecorder()

			handler.Login(w, req)
			if w.Code != tc.expectedStatus {
				t.Errorf("expected status code %d, got %d", tc.expectedStatus, w.Code)
			}

			if tc.checkResponse != nil {
				tc.checkResponse(t, w)
			}
		})
	}

}

func TestHandler_Logout(t *testing.T) {
	t.Parallel()

	cookieName := "rt"

	cases := []struct {
		name           string
		cookie         *http.Cookie
		busSetup       func(s *mockBus)
		expectedStatus int
		checkResponse  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "success: logs out user (with cookie)",
			cookie: &http.Cookie{
				Name:  cookieName,
				Value: "valid-refresh-token",
			},
			busSetup: func(b *mockBus) {
				b.On(Logout, nil)
			},
			expectedStatus: http.StatusNoContent,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				cookies := w.Result().Cookies()
				if len(cookies) != 1 {
					t.Errorf("expected 1 cookie, got %d", len(cookies))
				}
				if cookies[0].Name != cookieName {
					t.Errorf("expected cookie name to be %s, got %s", cookieName, cookies[0].Name)
				}
				if cookies[0].MaxAge != -1 {
					t.Errorf("expected cookie max age to be -1, got %d", cookies[0].MaxAge)
				}
			},
		},
		{
			name:           "success: logs out user (no cookie)",
			cookie:         nil,
			expectedStatus: http.StatusNoContent,
		},
		{
			name: "error: logout failed",
			cookie: &http.Cookie{
				Name:  cookieName,
				Value: "valid-refresh-token",
			},
			busSetup: func(b *mockBus) {
				b.On(Logout, business.ErrInternal)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			bus := newMockBus(t)
			if tc.busSetup != nil {
				tc.busSetup(bus)
			}

			handler := authhttp.New(bus, cookieName, "critiquefi.com")
			req := testutil.NewTestRequest(t, http.MethodPost, "/logout", nil)
			if tc.cookie != nil {
				req.AddCookie(tc.cookie)
			}
			w := httptest.NewRecorder()

			handler.Logout(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status code %d, got %d", tc.expectedStatus, w.Code)
			}

			if tc.checkResponse != nil {
				tc.checkResponse(t, w)
			}
		})
	}
}

func TestHandler_Refresh(t *testing.T) {
	t.Parallel()

	cookieName := "rt"

	cases := []struct {
		name           string
		cookie         *http.Cookie
		busSetup       func(s *mockBus)
		expectedStatus int
		checkResponse  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "success: refreshes token",
			cookie: &http.Cookie{
				Name:  cookieName,
				Value: "valid-refresh-token",
			},
			busSetup: func(b *mockBus) {
				b.On(Refresh, "access-token", "rotated-refresh-token", nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response authhttp.AuthenticationResponse
				testutil.DecodeResponse(t, w, &response)
				if response.AccessToken != "access-token" {
					t.Errorf("expected access token to be 'access-token', got %s", response.AccessToken)
				}

				cookies := w.Result().Cookies()
				if len(cookies) != 1 {
					t.Errorf("expected 1 cookie, got %d", len(cookies))
				}
				if cookies[0].Name != cookieName {
					t.Errorf("expected cookie name to be %s, got %s", cookieName, cookies[0].Name)
				}
				if cookies[0].Value != "rotated-refresh-token" {
					t.Errorf("expected cookie value to be 'rotated-refresh-token', got %s", cookies[0].Value)
				}
			},
		},
		{
			name:           "error: no token",
			cookie:         nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "error: failed to refresh token",
			cookie: &http.Cookie{
				Name:  cookieName,
				Value: "valid-refresh-token",
			},
			busSetup: func(b *mockBus) {
				b.On(Refresh, "", "", business.ErrInternal)
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			bus := newMockBus(t)
			if tc.busSetup != nil {
				tc.busSetup(bus)
			}

			handler := authhttp.New(bus, cookieName, "critiquefi.com")
			req := testutil.NewTestRequest(t, http.MethodPost, "/refresh", nil)
			if tc.cookie != nil {
				req.AddCookie(tc.cookie)
			}
			w := httptest.NewRecorder()

			handler.Refresh(w, req)
			if w.Code != tc.expectedStatus {
				t.Errorf("expected status code %d, got %d", tc.expectedStatus, w.Code)
			}
		})
	}
}
