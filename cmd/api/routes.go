package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/healthcheck", app.healthcheckHandler)

	mux.HandleFunc("GET /api/v1/chats", app.requireAuthenticatedUser(app.listChatHandler))
	mux.HandleFunc("GET /api/v1/chats/{uuid}", app.requireAuthenticatedUser(app.getChatHandler))
	mux.HandleFunc("GET /api/v1/chats/{uuid}/members", app.requireAuthenticatedUser(app.getChatMembers))

	mux.HandleFunc("POST /api/v1/users", app.createUserHandler)
	mux.HandleFunc("GET /api/v1/users/{uuid}", app.requireAuthenticatedUser(app.getUserHandler))

	mux.HandleFunc("POST /api/v1/tokens/auth", app.createAuthenticationTokenHandler)

	return app.authenticate(mux)
}
