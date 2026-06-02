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
