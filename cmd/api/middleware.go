package main

import (
	"fmt"
	"net/http"

	"golang.org/x/time/rate"
)

// recoverPanic is a middleware that recovers from a panic occurred in the same goroutine.
// Instead of an empty reply from the server, we defer a function to send a proper server Error to the user.
// The same strategy can be used if other goroutines are used to run some background tasks (not only HTTP requests).
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// rateLimit is a middleware that limits the number of requests received in a specific period of time.
// It implements a bucket rate limiter.
func (app *application) rateLimit(next http.Handler) http.Handler {
	// Initialize a new rate limiter which allows an average of 2 requests per second,
	// with a maximum of 4 requests in a single ‘burst’.
	// This initialization will run only once, when we wrap an Handler with this middleware.
	limiter := rate.NewLimiter(2, 4)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Call limiter.Allow() to see if the request is permitted, and if it's not,
		// then we call the rateLimitExceededResponse() helper to return a 429 Too Many
		// Requests response.
		if !limiter.Allow() {
			app.rateLimitExceededResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
