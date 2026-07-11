package testutil

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// NewTestRequest creates an HTTP request for testing
func NewTestRequest(t *testing.T, method, path string, body any) *http.Request {
	t.Helper()

	var req *http.Request
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		req = httptest.NewRequest(method, path, bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	return req
}

// DecodeResponse decodes the response body into the given struct
func DecodeResponse[T any](t *testing.T, w *httptest.ResponseRecorder, v *T) {
	t.Helper()
	err := json.NewDecoder(w.Body).Decode(v)
	if err != nil {
		t.Fatal(err)
	}
}
