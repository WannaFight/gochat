package main

import (
	"net/http"

	"github.com/WannaFight/gochat/ui"
	"github.com/a-h/templ"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.FS(ui.Files))

	mux.HandleFunc("GET /api/v1/healthcheck", app.healthcheckHandler)

	mux.HandleFunc("GET /api/v1/messages", app.requireAuthenticatedUser(app.getMessages))
	mux.HandleFunc("POST /api/v1/messages", app.requireAuthenticatedUser(app.createMessage))

	mux.HandleFunc("POST /api/v1/users", app.createUserHandler)
	mux.HandleFunc("GET /api/v1/users/{uuid}", app.requireAuthenticatedUser(app.getUserHandler))

	mux.HandleFunc("POST /api/v1/tokens/auth", app.createAuthenticationTokenHandler)

	mux.Handle("/login", templ.Handler(ui.LoginForm()))
	mux.Handle("/register", templ.Handler(ui.RegisterForm()))

	mux.Handle("/messages", templ.Handler(ui.MessageList()))

	mux.Handle("/assets/", fileServer)

	return app.AuthenticateWithCookie(mux)
}
