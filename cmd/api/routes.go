package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/healthcheck", app.healthcheckHandler)

	mux.HandleFunc("GET /api/v1/chats", app.listChatHandler)
	mux.HandleFunc("GET /api/v1/chats/{uuid}", app.getChatHandler)
	mux.HandleFunc("GET /api/v1/chats/{uuid}/members", app.getChatMembers)

	mux.HandleFunc("POST /api/v1/users", app.createUserHandler)
	mux.HandleFunc("GET /api/v1/users/{uuid}", app.getUserHandler)

	return mux
}
