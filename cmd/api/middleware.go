package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/WannaFight/gochat/internal/data"
)

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authHeader := r.Header.Get("Authorization")
		// If no Authorization header is provided, we can assume current user is anonymous.
		if authHeader == "" {
			r := app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Token" {
			app.logger.Debug("Authorization header is not valid", "header", authHeader)
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headerParts[1]
		if token == "" || len(token) > 26 {
			app.logger.Debug("Authorization token is not valid", "token", authHeader)
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		user, err := app.models.Users.GetByToken(token, data.ScopeAuthentication)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		r = app.contextSetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

// requireAuthenticatedUser check if user is presented in http.Request context.
// If it is missing or AnonymousUser return 401
func (app *application) requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)
		if user.IsAnonymous() {
			app.authenticationRequiredResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
