package utils

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// RespondWithJSON sends a JSON response with the given status code and payload
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error marshalling JSON response"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// RespondWithError sends an error response with the given status code and message
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

// ParseInt parses a string to an integer with error handling
func ParseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// ParseFloat parses a string to a float with error handling
func ParseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// ParseBool parses a string to a boolean with error handling
func ParseBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}
