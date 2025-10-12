package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, status int, message string) error {
	return WriteJSON(w, status, map[string]interface{}{
		"success": false,
		"error":   message,
	})
}

func WriteValidationError(w http.ResponseWriter, fields map[string]string) error {
	return WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
		"success": false,
		"error":   "Validation failed",
		"fields":  fields,
	})
}

func ParseJSON(r *http.Request, dest interface{}) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return fmt.Errorf("Content-Type must be application/json")
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(dest)
}
