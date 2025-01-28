package main

import (
	"context"
	"net/http"

	"github.com/AlessioPani/go-greenlight/internal/data"
)

// Context key custom type.
type contextKey string

// Key used for setting and getting the user from the request context.
const userContextKey = contextKey("user")

// contextSetUser sets the user into the request context.
func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(context.Background(), userContextKey, user)
	return r.WithContext(ctx)
}

// contextGetUser gets the user from the request context.
func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}
