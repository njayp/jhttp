package jhttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testStruct struct {
	Name string `json:"name"`
}

func (t testStruct) Valid() error {
	if t.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

func TestEncode(t *testing.T) {
	w := httptest.NewRecorder()
	v := testStruct{Name: "test"}

	err := Encode(w, http.StatusOK, v)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %v, got %v", http.StatusOK, w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("expected content type %q, got %q", "application/json", w.Header().Get("Content-Type"))
	}

	var result testStruct
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("expected no error decoding response, got %v", err)
	}

	if result.Name != v.Name {
		t.Fatalf("expected name %v, got %v", v.Name, result.Name)
	}
}

func TestEncodeError(t *testing.T) {
	w := httptest.NewRecorder()
	v := make(chan int)

	err := Encode(w, http.StatusOK, v)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %v, got %v", http.StatusInternalServerError, w.Code)
	}

	expected := "encode json: json: unsupported type: chan int\n"
	if w.Body.String() != expected {
		t.Fatalf("expected error message %q, got %q", expected, w.Body.String())
	}
}

func TestDecode(t *testing.T) {
	v := testStruct{Name: "test"}
	body, _ := json.Marshal(v)
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))

	w := httptest.NewRecorder()
	result, err := Decode[testStruct](w, r)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Name != v.Name {
		t.Fatalf("expected name %v, got %v", v.Name, result.Name)
	}
}

func TestDecodeInvalid(t *testing.T) {
	w := httptest.NewRecorder()
	v := testStruct{}
	body, _ := json.Marshal(v)
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))

	_, err := Decode[testStruct](w, r)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %v, got %v", http.StatusBadRequest, w.Code)
	}

	expected := "validation: name is required\n"
	if w.Body.String() != expected {
		t.Fatalf("expected error message %q, got %q", expected, w.Body.String())
	}
}
