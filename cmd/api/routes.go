package main

import (
	"net/http"

	"github.com/WannaFight/gochat/cmd/web"
	"github.com/a-h/templ"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.FS(web.Files))

	mux.HandleFunc("GET /api/v1/healthcheck", app.healthcheckHandler)

	mux.HandleFunc("GET /api/v1/chats", app.requireAuthenticatedUser(app.listChatHandler))
	mux.HandleFunc("POST /api/v1/chats", app.requireAuthenticatedUser(app.createChatHandler))
	mux.HandleFunc("GET /api/v1/chats/{uuid}", app.requireAuthenticatedUser(app.getChatHandler))
	mux.HandleFunc("GET /api/v1/chats/{uuid}/members", app.requireAuthenticatedUser(app.getChatMembers))
	mux.HandleFunc("POST /api/v1/chats/{uuid}/members", app.requireAuthenticatedUser(app.createChatMember))
	mux.HandleFunc("GET /api/v1/chats/{uuid}/messages", app.requireAuthenticatedUser(app.getChatMessages))
	mux.HandleFunc("POST /api/v1/chats/{uuid}/messages", app.requireAuthenticatedUser(app.createChatMessage))

	mux.HandleFunc("POST /api/v1/users", app.createUserHandler)
	mux.HandleFunc("GET /api/v1/users/{uuid}", app.requireAuthenticatedUser(app.getUserHandler))

	mux.HandleFunc("POST /api/v1/tokens/auth", app.createAuthenticationTokenHandler)

	mux.Handle("/login", templ.Handler(web.LoginForm()))
	mux.Handle("/register", templ.Handler(web.RegisterForm()))
	mux.Handle("/assets/", fileServer)

	return app.authenticate(mux)
}
