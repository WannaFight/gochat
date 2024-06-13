package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/WannaFight/gochat/internal/data"
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
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{Username: input.Username}
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateUsername):
			app.badRequestResponse(w, r, errors.New("user with this username already exists"))
		default:
			app.serverErrorResponse(w, r)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r)
	}
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
