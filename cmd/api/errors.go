package main

import "net/http"

func (app *application) errorResponse(w http.ResponseWriter, status int, message any) {
	env := envelope{"error": message}

	if err := app.writeJSON(w, status, env, nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request) {
	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, http.StatusInternalServerError, message)
}
