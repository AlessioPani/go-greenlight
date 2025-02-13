package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlessioPani/go-greenlight/internal/data/mocks"
)

// Test method used to test the createAuthenticationTokenHandler http.Handler.
func TestCreateAuthenticationTokenHandler(t *testing.T) {
	// Get test application config and handler.
	app := newTestApplication()
	handler := http.Handler(app.routes())

	// Tests to be run.
	var tests = []struct {
		name           string
		method         string
		url            string
		email          string
		passwords      string
		expectedResult int
	}{
		{"valid token", "POST", "/v1/tokens/authentication", mocks.ActiveUser.Email, "valid_password", http.StatusOK},
		{"invalid email", "POST", "/v1/tokens/authentication", "mail", "valid_password", http.StatusUnprocessableEntity},
		{"invalid credentials", "POST", "/v1/tokens/authentication", "invalid@mail.com", "valid_password", http.StatusUnauthorized},
	}

	// Executes tests.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create the handler input.
			var input struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}
			input.Email = test.email
			input.Password = test.passwords

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

// Test method used to test the createPasswordResetToken http.Handler.
func TestCreatePasswordResetTokenHandler(t *testing.T) {
	// Get test application config and handler.
	app := newTestApplication()
	handler := http.Handler(app.routes())

	// Tests to be run.
	var tests = []struct {
		name           string
		method         string
		url            string
		email          string
		expectedResult int
	}{
		{"valid token", "POST", "/v1/tokens/password-reset", mocks.ActiveUser.Email, http.StatusAccepted},
		{"invalid email", "POST", "/v1/tokens/password-reset", "mail", http.StatusUnprocessableEntity},
		{"invalid credentials", "POST", "/v1/tokens/password-reset", "invalid@mail.com", http.StatusUnprocessableEntity},
	}

	// Executes tests.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create the handler input.
			var input struct {
				Email string `json:"email"`
			}
			input.Email = test.email

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

// Test method used to test the createActivationTokenHandler http.Handler.
func TestCreateActivationTokenHandler(t *testing.T) {
	// Get test application config and handler.
	app := newTestApplication()
	handler := http.Handler(app.routes())

	// Tests to be run.
	var tests = []struct {
		name           string
		method         string
		url            string
		email          string
		expectedResult int
	}{
		{"valid user", "POST", "/v1/tokens/activation", mocks.InactiveUser.Email, http.StatusAccepted},
		{"invalid user", "POST", "/v1/tokens/activation", mocks.ActiveUser.Email, http.StatusUnprocessableEntity},
		{"invalid email", "POST", "/v1/tokens/activation", "mail", http.StatusUnprocessableEntity},
		{"invalid credentials", "POST", "/v1/tokens/activation", "invalid@mail.com", http.StatusUnprocessableEntity},
	}

	// Executes tests.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create the handler input.
			var input struct {
				Email string `json:"email"`
			}
			input.Email = test.email

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
