package main

import (
	"net/http"
	"testing"
)

// Test method used to test the routes function.
func TestRoutes(t *testing.T) {
	app := newTestApplication()

	router := app.routes()

	switch rt := router.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Errorf("got %s type, expected %s type", rt, "http.Handler")
	}
}
