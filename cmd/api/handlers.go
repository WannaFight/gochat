package main

import (
	"fmt"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r)
	}
}

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"username"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		// TODO: error handling
		return
	}
	fmt.Fprintln(w, "creating user")
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	uuid, err := app.readUUIDParam(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err.Error())
		return
	}
	fmt.Fprintf(w, "getting user with uuid=%s", uuid.String())
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

func (app *application) getChatMembers(w http.ResponseWriter, r *http.Request) {
	uuid, err := app.readUUIDParam(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err.Error())
		return
	}
	fmt.Fprintf(w, "listing all chat memebers of chat with uuid=%s\n", uuid.String())
}
