package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

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
// It implements an ip-based rate limiter, using mutex to safely add entries in the clients map
// and a goroutine to delete the older entries (in order to not have a map growning indefinitely).
func (app *application) rateLimit(next http.Handler) http.Handler {
	// client is a struct used to hold the rate limiter and the last seen time for each
	// client.
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	// Declare a mutex and a map to hold the clients's IP addresses and rate limiters.
	var mu sync.Mutex
	var clients = make(map[string]*client)

	// Launch a goroutine which removes old entries from the clients map once every minute.
	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the client's IP address from the request.
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		// Maps are not safe for concurrent use. We use mutex to safely interact with the clients map.
		mu.Lock()

		// For the current client, set up a new limiter of 2 requests per second with a maximum
		// of 4 requests in a single burst.
		if _, found := clients[ip]; !found {
			clients[ip] = &client{limiter: rate.NewLimiter(2, 4)}
		}

		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			app.rateLimitExceededResponse(w, r)
			return
		}

		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
