package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Test method used to test the healthcheck endpoint.
func TestHealthcheckHandler(t *testing.T) {
	// Get test application config and handler.
	app := newTestApplication()
	handler := http.Handler(app.routes())

	// Tests to be run.
	tests := []struct {
		name           string
		method         string
		target         string
		expectedStatus int
		expectedBody   string
	}{
		{"ValidRequest", "GET", "/v1/healthcheck", http.StatusOK, "available"},
	}

	// Executes tests.
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test request.
			req, err := http.NewRequest(tc.method, tc.target, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create a test response recorder (response writer).
			rw := httptest.NewRecorder()

			// Serve test request.
			handler.ServeHTTP(rw, req)

			// Check for test results.
			if status := rw.Code; status != tc.expectedStatus {
				t.Errorf("expected status %v, got %v", tc.expectedStatus, status)
			}

			if !strings.Contains(rw.Body.String(), tc.expectedBody) {
				t.Errorf("expected body %v, got %v", tc.expectedBody, rw.Body.String())
			}
		})
	}
}
