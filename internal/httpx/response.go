package httpx

import (
	"encoding/json"
	"net/http"
)

// WriteJSON serializes payload as JSON and writes it to the response writer.
func WriteJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if payload == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(payload)
}

// WriteError is a helper to return consistent JSON error payloads.
func WriteError(w http.ResponseWriter, statusCode int, message string) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	WriteJSON(w, statusCode, errorResponse{Error: message})
}
