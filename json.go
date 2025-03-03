package jhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Encode writes a value of type T to the writer as a JSON response with the
// given status code.
func Encode[T any](w http.ResponseWriter, status int, v T) error {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		err = fmt.Errorf("encode: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return nil
}

type valid interface {
	Valid() error
}

// Decode decodes the request body into a value of type T and validates it. If
// the value is not valid, it writes a bad request response to the writer and
// returns an error.
func Decode[T valid](w http.ResponseWriter, r *http.Request) (T, error) {
	v, err := decode[T](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	return v, err
}

func decode[T valid](r *http.Request) (T, error) {
	var v T

	// decode the request body into v
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode: %w", err)
	}

	// validate the decoded value
	if err := v.Valid(); err != nil {
		return v, fmt.Errorf("validation: %w", err)
	}

	return v, nil
}
