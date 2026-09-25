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

// TestCreateMovieHandler covers successful creation and validation or model failures.
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

// TestShowMovieHandler covers found, missing, malformed, and unsupported requests.
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

// TestUpdateMovieHandler covers full updates, validation, missing records, and conflicts.
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

// TestUpdateMovieHandlerPartialAndInvalidPayloads covers sparse patches and rejected values.
func TestUpdateMovieHandlerPartialAndInvalidPayloads(t *testing.T) {
	tests := []struct {
		name            string
		body            string
		expectedVersion string
		wantStatus      int
	}{
		{"no fields", `{}`, "1", http.StatusOK},
		{"null field", `{"title":null}`, "1", http.StatusBadRequest},
		{"wrong field type", `{"year":"not a year"}`, "1", http.StatusBadRequest},
		{"version mismatch", `{}`, "2", http.StatusConflict},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPatch, "/v1/movies/1", bytes.NewBufferString(tt.body))
			req = req.WithContext(context.WithValue(req.Context(), userContextKey, mocks.ActiveUser))
			req.Header.Set("Authorization", "Bearer 12345678901234567890123456")
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Expected-Version", tt.expectedVersion)
			req.RemoteAddr = "localhost:8080"
			rr := httptest.NewRecorder()
			newTestApplication().routes().ServeHTTP(rr, req)
			if rr.Code != tt.wantStatus {
				t.Fatalf("status %d, want %d; body %s", rr.Code, tt.wantStatus, rr.Body.String())
			}
		})
	}
}

// TestDeleteMovieHandler covers successful deletion and missing records.
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

// TestListMovieHandler covers valid filters, invalid filters, and model failures.
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
