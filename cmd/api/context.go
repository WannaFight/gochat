package main

import (
	"context"
	"net/http"

	"github.com/WannaFight/gochat/internal/data"
)

type contextKey string

const (
	contextUserKey = contextKey("user")
)

// contextSetUser stores passed user to http.Request context.
func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), contextUserKey, user)
	return r.WithContext(ctx)
}

// contextGetUser retrieves user from http.Request context.
func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(contextUserKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}
