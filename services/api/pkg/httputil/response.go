package httputil

import (
	"encoding/json"
	"net/http"
)

// WriteJSON sets the HTTP response data
func WriteJSON(w http.ResponseWriter, status int, data any) {
	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(b); err != nil {
		return
	}
}

// WriteError writes a generic error response
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"error": message})
}

// WriteBadRequest writes a 400 bad request response
func WriteBadRequest(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusBadRequest, message)
}

// WriteUnauthorized writes a 401 unauthorized response
func WriteUnauthorized(w http.ResponseWriter) {
	WriteError(w, http.StatusUnauthorized, "unauthorized")
}

// WriteForbidden writes a 403 forbidden response
func WriteForbidden(w http.ResponseWriter) {
	WriteError(w, http.StatusForbidden, "forbidden")
}

// WriteNotFound writes a 404 not found response
func WriteNotFound(w http.ResponseWriter) {
	WriteError(w, http.StatusNotFound, "not found")
}

// WriteConflict writes a 409 conflict response
func WriteConflict(w http.ResponseWriter) {
	WriteError(w, http.StatusConflict, "conflict")
}

// WriteUnprocessable writes a 422 unprocessable entity response
func WriteUnprocessable(w http.ResponseWriter, errors any) {
	WriteJSON(w, http.StatusUnprocessableEntity, errors)
}

// WriteInternalError writes a 500 internal server error response
func WriteInternalError(w http.ResponseWriter) {
	WriteError(w, http.StatusInternalServerError, "internal server error")
}

// WriteServiceUnavailable writes a 503 service unavailable response
func WriteServiceUnavailable(w http.ResponseWriter) {
	WriteError(w, http.StatusServiceUnavailable, "service unavailable")
}

// WriteNotImplemented writes a 501 not implemented response
func WriteNotImplemented(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotImplemented)
	if _, err := w.Write([]byte("not implemented")); err != nil {
		return
	}
}

// WriteNoContent writes a 204 no content response
func WriteNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
