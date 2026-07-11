// Package testutil provides shared test helpers.
// This should ONLY be consumed by test packages.
package testutil

import (
	"errors"
	"testing"

	"github.com/andrew-hayworth22/critiquefi-service/internal/models"
)

// CheckErr compares the actual error to the expected error.
// Includes support for validation errors.
func CheckErr(actual error, expected error, t *testing.T) {
	t.Helper()

	if expected == nil {
		if actual != nil {
			t.Fatalf("expected no error, got %v", actual)
		}
		return
	}

	var expectedVE models.ValidationErrors
	if errors.As(expected, &expectedVE) {
		var actualVE models.ValidationErrors
		if !errors.As(actual, &actualVE) {
			t.Fatalf("expected validation errors, got %v", actual)
		}
		if !expectedVE.Equals(actualVE) {
			t.Fatalf("expected validation errors %v, got %v", expectedVE, actualVE)
		}
		return
	}

	if !errors.Is(actual, expected) {
		t.Fatalf("expected error %v, got %v", expected, actual)
	}
}

// ConvertError safely converts an any value to an error
func ConvertError(v any) error {
	if v == nil {
		return nil
	}
	return v.(error)
}
