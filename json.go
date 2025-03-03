package jhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Encode[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}

	return nil
}

type valid interface {
	Valid() error
}

func Decode[T valid](r *http.Request) (T, error) {
	var v T

	// decode the request body into v
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}

	// validate the decoded value
	if err := v.Valid(); err != nil {
		return v, fmt.Errorf("validate: %w", err)
	}

	return v, nil
}
