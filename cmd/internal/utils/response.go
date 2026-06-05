package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

// ApiResponse writes a JSON response with given status code.
func ApiResponse(w http.ResponseWriter, statusCode int, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("failed to encode json response:", err)
	}
}

// DecodeJSONBody decodes a JSON request body into the provided struct.
func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Disallow unknown fields to prevent silent errors

	if err := decoder.Decode(dst); err != nil {
		return err
	}

	return nil
}
