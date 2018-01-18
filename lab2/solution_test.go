// CMPT 315 (Winter 2018)
// Lab #2: Introduction to gorilla/mux
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/gorilla/mux"
)

func TestGet(t *testing.T) {
	// TODO: define the *mux.Router in one place
	// TODO: router := routes()
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/digits/{id}", handleGetDigit)

	tests := []struct {
		url        string
		statusCode int
		body       string
	}{
		{
			url:        "/api/v1/digits/0",
			statusCode: http.StatusOK,
			body:       "zero",
		},
		{
			url:        "/api/v1/digits/9",
			statusCode: http.StatusOK,
			body:       "nine",
		},
		{
			url:        "/api/v1/digits/100",
			statusCode: http.StatusGone,
			body:       "",
		},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodGet, tt.url, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		resp := w.Result()

		if resp.StatusCode != tt.statusCode {
			t.Errorf("Expected %d, got %d for %v", tt.statusCode, resp.StatusCode, tt.url)
		}

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		body := strings.TrimSpace(string(b))
		if body != tt.body {
			t.Errorf("Expected %q, got %q for %v", tt.body, body, tt.url)
		}
	}
}

func TestDelete(t *testing.T) {
	// TODO: router := routes()
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/digits/{id}", handleDeleteDigit)

	tests := []struct {
		url        string
		statusCode int
	}{
		{
			url:        "/api/v1/digits/9",
			statusCode: http.StatusNoContent,
		},
		{
			url:        "/api/v1/digits/100",
			statusCode: http.StatusGone,
		},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodDelete, tt.url, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		resp := w.Result()

		if resp.StatusCode != tt.statusCode {
			t.Errorf("Expected %d, got %d for %v", tt.statusCode, resp.StatusCode, tt.url)
		}
	}
}

func TestGetAfterDelete(t *testing.T) {
	// TODO: router := routes()
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/digits/{id:[0-9]}", handleGetDigit).Methods("GET")
	router.HandleFunc("/api/v1/digits/{id:[0-9]}", handleDeleteDigit).Methods("DELETE")

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/digits/0", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("DELETE expected %d, got %d", http.StatusNoContent, resp.StatusCode)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/v1/digits/0", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusGone {
		t.Fatalf("GET expected %d, got %d", http.StatusGone, resp.StatusCode)
	}
}

// NOTE: run with the race detector: go test -race
func TestConcurrentDelete(t *testing.T) {
	// TODO: router := routes()
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/digits/{id:[0-9]}", handleDeleteDigit).Methods("DELETE")

	var wg sync.WaitGroup

	for d := 1; d <= 8; d++ {
		wg.Add(1)
		go func(d int) {
			url := fmt.Sprintf("/api/v1/digits/%d", d)
			req := httptest.NewRequest(http.MethodDelete, url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			resp := w.Result()
			if resp.StatusCode != http.StatusNoContent {
				t.Errorf("DELETE expected %d, got %d", http.StatusNoContent, resp.StatusCode)
			}
			wg.Done()
		}(d)
	}
	wg.Wait()
}

// TODO: Add additional tests for invalid routes (404) and non-numeric parameters.

// TODO: A helper function to reset the digits map before each test would make these tests more
// independent and resilient to change. Make this pass: go test -count 3
