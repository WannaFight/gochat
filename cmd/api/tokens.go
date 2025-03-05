package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/WannaFight/gochat/internal/data"
)

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.models.Users.GetByUsername(input.Username)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.badRequestResponse(w, r, errors.New("invalid credentials"))
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	isPasswordCorrect, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !isPasswordCorrect {
		app.badRequestResponse(w, r, errors.New("invalid credentials"))
		return
	}

	// Do not store old tokens. TODO: maybe move to goroutine
	err = app.models.Tokens.DeleteForUser(user.ID, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	token, err := app.models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := app.generateTokenCookie(token.PlainText)
	headers.Add("HX-Redirect", "/chats")

	err = app.writeJSON(w, http.StatusCreated, envelope{"token": token}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
