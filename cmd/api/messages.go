package main

import (
	"net/http"

	"github.com/WannaFight/gochat/internal/data"
)

func (app *application) getMessages(w http.ResponseWriter, r *http.Request) {
	messages, err := app.models.Messages.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"messages": messages}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createMessage(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Text string `json:"text"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	message := &data.Message{
		Text:   input.Text,
		Author: *app.contextGetUser(r),
	}

	if err = app.models.Messages.Insert(message); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"message": message}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
