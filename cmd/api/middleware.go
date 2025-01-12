package main

import (
	"fmt"
	"net/http"
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
