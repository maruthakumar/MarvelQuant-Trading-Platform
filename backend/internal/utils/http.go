package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// RespondWithError sends an error response with the specified status code and message
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

// RespondWithJSON sends a JSON response with the specified status code and payload
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// Convert payload to JSON
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Error marshalling JSON response"}`))
		return
	}

	// Set headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// ParsePaginationParams parses pagination parameters from request
func ParsePaginationParams(r *http.Request) (page, limit int) {
	// Default values
	page = 1
	limit = 20

	// Parse query parameters
	query := r.URL.Query()
	
	// Parse page
	if pageStr := query.Get("page"); pageStr != "" {
		if parsedPage, err := parseInt(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}
	
	// Parse limit
	if limitStr := query.Get("limit"); limitStr != "" {
		if parsedLimit, err := parseInt(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}
	
	return page, limit
}

// parseInt parses a string to an integer
func parseInt(s string) (int, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}
