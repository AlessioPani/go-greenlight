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

// Test method used for the registerUserHandler http.Handler.
func TestRegisterUserHandler(t *testing.T) {
	// Get test application config and handler.
	app := newTestApplication()
	handler := http.Handler(app.routes())

	// Tests to be run.
	tests := []struct {
		name           string
		method         string
		url            string
		payload        data.User
		password       string
		expectedResult int
	}{
		{"user created", "POST", "/v1/users", mocks.InactiveUser, "test_password", http.StatusCreated},
		{"email already used", "POST", "/v1/users", mocks.ActiveUser, "test_password", http.StatusUnprocessableEntity},
		{"invalid user", "POST", "/v1/users", mocks.InvalidUser, "test_password", http.StatusUnprocessableEntity},
	}

	// Execute tests.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create the handler input.
			var input struct {
				Name     string `json:"name"`
				Email    string `json:"email"`
				Password string `json:"password"`
			}
			input.Name = test.payload.Name
			input.Email = test.payload.Email
			input.Password = test.password

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

// Test method used for the activateUserHandler http.Handler.
func TestActivateUserHandler(t *testing.T) {
	{
		// Get test application config and handler.
		app := newTestApplication()
		handler := http.Handler(app.routes())

		// Tests to be run.
		tests := []struct {
			name           string
			method         string
			url            string
			token          string
			expectedResult int
		}{
			{"valid token", "PUT", "/v1/users/activated", "validtoken1234567890token0", http.StatusOK},
			{"invalid token", "PUT", "/v1/users/activated", "", http.StatusUnprocessableEntity},
			{"expired token", "PUT", "/v1/users/activated", "expiredtoken123456789token", http.StatusUnprocessableEntity},
		}

		// Execute tests.
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				// Create the handler input.
				var input struct {
					TokenPlaintext string `json:"token"`
				}
				input.TokenPlaintext = test.token

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
}

// Test method used for the updateUserPasswordHandler http.handler
func TestUpdateUserPasswordHandler(t *testing.T) {
	// Get test application config and handler.
	app := newTestApplication()
	handler := http.Handler(app.routes())

	// Tests to be run.
	tests := []struct {
		name           string
		method         string
		url            string
		password       string
		token          string
		expectedResult int
	}{
		{"valid password", "PUT", "/v1/users/password", "valid_password", "validtoken1234567890token0", http.StatusOK},
		{"invalid password", "PUT", "/v1/users/password", "", "", http.StatusUnprocessableEntity},
	}

	// Execute tests.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create the handler input.
			var input struct {
				Password       string `json:"password"`
				TokenPlaintext string `json:"token"`
			}
			input.Password = test.password
			input.TokenPlaintext = test.token

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
