package main

import (
	"fmt"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "status: available")
	fmt.Fprintf(w, "environment: %s\n", app.config.env)
	fmt.Fprintf(w, "version: %s\n", app.config.version)
}

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "creating user")
}

func (app *application) listChatHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "listing all user chats")
}

func (app *application) getChatHandler(w http.ResponseWriter, r *http.Request) {
	uuid, err := app.readUUIDParam(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err.Error())
		return
	}
	fmt.Fprintf(w, "getting chat with uuid=%s\n", uuid.String())
}
