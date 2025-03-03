package jhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Encode writes a value of type T to the response body.
func Encode[T any](w http.ResponseWriter, status int, v T) error {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return nil
}

// Decode reads a value of type T from the request body.
func Decode[T any](r *http.Request) (T, error) {
	var v T

	// decode the request body into v
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode: %w", err)
	}

	return v, nil
}
