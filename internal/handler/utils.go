// Package handler provides HTTP request handlers for the post-comments REST API.
// This package contains the HTTP layer logic, including request/response handling,
// JSON serialization, error handling, and business logic coordination.
package handler

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents a standardized error response format for the API.
// This struct ensures consistent error messaging across all endpoints.
type ErrorResponse struct {
	Error string `json:"error"`
}

// RespondWithError sends a standardized error response with the specified HTTP status code.
// This function ensures consistent error response format across all API endpoints.
//
// Parameters:
//   - w: HTTP response writer
//   - code: HTTP status code (e.g., 400, 404, 500)
//   - message: Human-readable error message
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, ErrorResponse{Error: message})
}

// RespondWithJSON sends a JSON response with the specified HTTP status code.
// This function handles JSON marshaling and sets appropriate headers.
//
// Parameters:
//   - w: HTTP response writer
//   - code: HTTP status code (e.g., 200, 201, 400)
//   - payload: Any JSON-serializable data structure
//
// If JSON marshaling fails, sends a 500 Internal Server Error response.
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
