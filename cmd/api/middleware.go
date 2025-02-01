package main

import (
	"errors"
	"expvar"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/AlessioPani/go-greenlight/internal/data"
	"github.com/AlessioPani/go-greenlight/internal/validator"
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
		if app.config.limiter.enabled {
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
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst),
				}
			}

			clients[ip].lastSeen = time.Now()

			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.rateLimitExceededResponse(w, r)
				return
			}

			mu.Unlock()
		}

		next.ServeHTTP(w, r)
	})
}

// authenticate is a middleware that handle authentication by setting the request context.
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This indicates to any caches that the response may vary based on the value of the
		// Authorization header in the request.
		w.Header().Add("Vary", "Authorization")

		// Check for the authorization header. If it is empty put in the context the anonymous user.
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			r = app.contextSetUser(r, data.AnonymouseUser)
			next.ServeHTTP(w, r)
			return
		}

		// Check for the authorization header to be in the Bearer <token> format.
		// If not, send the appropriate error to the user.
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		// Gets the actual token and validates it.
		token := headerParts[1]

		v := validator.New()

		if data.ValidateTokenPlaintext(v, token); !v.Valid() {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		// Retrieve the user from the DB by its auth token.
		user, err := app.models.Users.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		// Puts the user in the context of the new request.
		r = app.contextSetUser(r, user)

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}

// requireAuthenticatedUser is a middleware that verifies if a user is authenticated
// and can get the requested resource.
// Accepts and returns a http.HandlerFunc to encapsulate the handlers in routes.go
// without any further conversion.
func (app *application) requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the user from the request context.
		user := app.contextGetUser(r)

		if user.IsAnonymous() {
			app.authenticationRequiredResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// requireActivatedUser is a middleware that verifies if a user is activated
// and can get the requested resource.
// Accepts and returns a http.HandlerFunc to encapsulate the handlers in routes.go
// without any further conversion.
func (app *application) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the user from the request context.
		user := app.contextGetUser(r)

		if !user.Activated {
			app.inactiveAccountResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})

	return app.requireAuthenticatedUser(fn)
}

// requirePermission is a middleware that verifies if a user has a required permission.
func (app *application) requirePermission(code string, next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Get the user from the request context.
		user := app.contextGetUser(r)

		// Get the slice of permissions for the user.
		permissions, err := app.models.Permissions.GetAllForUser(user.ID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		// Check if the slice includes the required permission. If it doesn't, then
		// return a 403 Forbidden response.
		if !permissions.Include(code) {
			app.notPermittedResponse(w, r)
			return
		}

		// Otherwise call the next handler in the chain.
		next.ServeHTTP(w, r)
	}

	return app.requireActivatedUser(fn)
}

// metrics is a middleware that records several request-level metrics.
func (app *application) metrics(next http.Handler) http.Handler {
	// Request-level metrics
	var (
		totalRequestsReceived           = expvar.NewInt("total_requests_received")
		totalResponsesSent              = expvar.NewInt("total_responses_sent")
		totalProcessingTimeMicroseconds = expvar.NewInt("total_processing_time_us")
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		totalRequestsReceived.Add(1)

		next.ServeHTTP(w, r)

		totalResponsesSent.Add(1)

		duration := time.Since(start).Microseconds()
		totalProcessingTimeMicroseconds.Add(duration)
	})
}
