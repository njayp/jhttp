package jhttp

import (
	"encoding/json"
	"net/http"
)

// Encode writes a value of type T to the response body.
func Encode[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// Decode reads a value of type T from the request body.
func Decode[T any](r *http.Request) (T, error) {
	var v T
	return v, json.NewDecoder(r.Body).Decode(&v)
}
