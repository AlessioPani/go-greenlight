package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlessioPani/go-greenlight/internal/data"
	"github.com/AlessioPani/go-greenlight/internal/data/mocks"
)

// Test method used for the createMovieHandler http.Handler.
func TestCreateMovieHandler(t *testing.T) {
	// Get test application config and handler.
	app := newTestApplication()
	handler := http.Handler(app.routes())

	// Tests to be run.
	tests := []struct {
		name           string
		method         string
		url            string
		payload        data.Movie
		expectedResult int
	}{
		{"movie created", "POST", "/v1/movies", mocks.ValidMovie, http.StatusCreated},
		{"movie invalid", "POST", "/v1/movies", mocks.InvalidMovie, http.StatusUnprocessableEntity},
		{"db error", "POST", "/v1/movies", mocks.ErrorMovie, http.StatusInternalServerError},
	}

	// Execute tests.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create the handler input.
			var input struct {
				Title   string       `json:"title"`
				Year    int32        `json:"year"`
				Runtime data.Runtime `json:"runtime"`
				Genres  []string     `json:"genres"`
			}
			input.Title = test.payload.Title
			input.Year = test.payload.Year
			input.Runtime = test.payload.Runtime
			input.Genres = test.payload.Genres

			// Create a test request with body, context and authorization.
			js, err := json.Marshal(input)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(test.method, test.url, bytes.NewBuffer(js))
			if err != nil {
				t.Fatal(err)
			}

			ctx := context.WithValue(context.Background(), userContextKey, mocks.ActiveUser)
			req = req.WithContext(ctx)
			req.Header.Set("Authorization", "Bearer 12345678901234567890123456")
			req.Header.Set("Content-Type", "application/json")
			req.RemoteAddr = "localhost:8080"

			// Create a test response recorder (response writer).
			rw := httptest.NewRecorder()

			// Serve test request.
			handler.ServeHTTP(rw, req)

			// Check test results.
			if rw.Code != test.expectedResult {
				t.Errorf("got %d, expected %d", rw.Code, test.expectedResult)
			}
		})
	}
}

// Test method used for the showMovieHandler http.Handler.
func TestShowMovieHandler(t *testing.T) {
	// Get test application config and handler.
	app := newTestApplication()
	handler := http.Handler(app.routes())

	// Tests to be run.
	tests := []struct {
		name           string
		method         string
		url            string
		expectedResult int
	}{
		{"movie found", "GET", "/v1/movies/1", http.StatusOK},
		{"movie not found", "GET", "/v1/movies/4", http.StatusNotFound},
		{"invalid url", "GET", "/v1/movies/foo", http.StatusNotFound},
		{"server error", "GET", "/v1/movies/5", http.StatusInternalServerError},
		{"invalid method", "POST", "/v1/movies/1", http.StatusMethodNotAllowed},
	}

	// Execute tests.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a test request with context and authorization.
			req, err := http.NewRequest(test.method, test.url, nil)
			if err != nil {
				t.Fatal(err)
			}
			ctx := context.WithValue(context.Background(), userContextKey, mocks.ActiveUser)
			req = req.WithContext(ctx)
			req.Header.Set("Authorization", "Bearer 12345678901234567890123456")
			req.RemoteAddr = "localhost:8080"

			// Create a test response recorder (response writer).
			rw := httptest.NewRecorder()

			// Serve test request.
			handler.ServeHTTP(rw, req)

			// Check test results.
			if rw.Code != test.expectedResult {
				t.Errorf("got %d, expected %d", rw.Code, test.expectedResult)
			}
		})
	}
}

// Test method used for the updateMovieHandler http.Handler.
func TestUpdateMovieHandler(t *testing.T) {
	// Get test application config and handler.
	app := newTestApplication()
	handler := http.Handler(app.routes())

	// Tests to be run.
	tests := []struct {
		name           string
		method         string
		url            string
		payload        data.Movie
		expectedResult int
	}{
		{"movie found", "PATCH", "/v1/movies/1", mocks.ValidMovie, http.StatusOK},
		{"invalid movie", "PATCH", "/v1/movies/2", mocks.InvalidMovie, http.StatusUnprocessableEntity},
		{"error movie", "PATCH", "/v1/movies/3", mocks.ErrorMovie, http.StatusConflict},
		{"movie not found", "PATCH", "/v1/movies/4", data.Movie{}, http.StatusNotFound},
		{"invalid url", "PATCH", "/v1/movies/foo", data.Movie{}, http.StatusNotFound},
	}

	// Execute tests.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create the handler input.
			var input struct {
				Title   *string       `json:"title"`
				Year    *int32        `json:"year"`
				Runtime *data.Runtime `json:"runtime"`
				Genres  []string      `json:"genres"`
			}
			input.Title = &test.payload.Title
			input.Year = &test.payload.Year
			input.Runtime = &test.payload.Runtime
			input.Genres = test.payload.Genres

			// Create a test request with body, context and authorization.
			js, err := json.Marshal(input)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(test.method, test.url, bytes.NewBuffer(js))
			if err != nil {
				t.Fatal(err)
			}

			ctx := context.WithValue(context.Background(), userContextKey, mocks.ActiveUser)
			req = req.WithContext(ctx)
			req.Header.Set("Authorization", "Bearer 12345678901234567890123456")
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Expected-Version", "1")
			req.RemoteAddr = "localhost:8080"

			// Create a test response recorder (response writer).
			rw := httptest.NewRecorder()

			// Serve test request.
			handler.ServeHTTP(rw, req)

			// Check test results.
			if rw.Code != test.expectedResult {
				t.Errorf("got %d, expected %d", rw.Code, test.expectedResult)
			}
		})
	}
}

// Test method used for the deleteMovieHandler http.Handler.
func TestDeleteMovieHandler(t *testing.T) {
	// Get test application config and handler.
	app := newTestApplication()
	handler := http.Handler(app.routes())

	// Tests to be run.
	tests := []struct {
		name           string
		method         string
		url            string
		expectedResult int
	}{
		{"movie found", "DELETE", "/v1/movies/1", http.StatusOK},
		{"invalid url", "DELETE", "/v1/movies/foo", http.StatusNotFound},
		{"movie not found", "DELETE", "/v1/movies/4", http.StatusNotFound},
	}

	// Execute tests.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a test request with context and authorization.
			req, err := http.NewRequest(test.method, test.url, nil)
			if err != nil {
				t.Fatal(err)
			}
			ctx := context.WithValue(context.Background(), userContextKey, mocks.ActiveUser)
			req = req.WithContext(ctx)
			req.Header.Set("Authorization", "Bearer 12345678901234567890123456")
			req.RemoteAddr = "localhost:8080"

			// Create a test response recorder (response writer).
			rw := httptest.NewRecorder()

			// Serve test request.
			handler.ServeHTTP(rw, req)

			// Check test results.
			if rw.Code != test.expectedResult {
				t.Errorf("got %d, expected %d", rw.Code, test.expectedResult)
			}
		})
	}

}

// Test method used for the listMovieHandler http.Handler.
func TestListMovieHandler(t *testing.T) {
	// Get test application config and handler.
	app := newTestApplication()
	handler := http.Handler(app.routes())

	// Tests to be run.
	tests := []struct {
		name           string
		method         string
		url            string
		expectedResult int
	}{
		{"valid request", "GET", "/v1/movies?title=rings", http.StatusOK},
		{"invalid request", "GET", "/v1/movies?sort=invalid", http.StatusUnprocessableEntity},
		{"server error", "GET", "/v1/movies?title=watchmen", http.StatusInternalServerError},
	}

	// Execute tests.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a test request with context and authorization.
			req, err := http.NewRequest(test.method, test.url, nil)
			if err != nil {
				t.Fatal(err)
			}
			ctx := context.WithValue(context.Background(), userContextKey, mocks.ActiveUser)
			req = req.WithContext(ctx)
			req.Header.Set("Authorization", "Bearer 12345678901234567890123456")
			req.RemoteAddr = "localhost:8080"

			// Create a test response recorder (response writer).
			rw := httptest.NewRecorder()

			// Serve test request.
			handler.ServeHTTP(rw, req)

			// Check test results.
			if rw.Code != test.expectedResult {
				t.Errorf("got %d, expected %d", rw.Code, test.expectedResult)
			}
		})
	}
}
