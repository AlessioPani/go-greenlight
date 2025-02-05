package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthcheckHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		target         string
		expectedStatus int
		expectedBody   string
	}{
		{"ValidRequest", "GET", "/v1/healthcheck", http.StatusOK, "available"},
	}

	app := newTestApplication()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, tc.target, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(app.healthcheckHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("expected status %v, got %v", tc.expectedStatus, status)
			}

			if !strings.Contains(rr.Body.String(), tc.expectedBody) {
				t.Errorf("expected body %v, got %v", tc.expectedBody, rr.Body.String())
			}
		})
	}
}
