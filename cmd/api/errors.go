package main

import "net/http"

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	if err := app.writeJSON(w, status, env, nil); err != nil {
		app.logger.Error("internal server error",
			"error", err.Error(),
			"request_method", r.Method,
			"request_url", r.URL.String(),
		)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	message := "the server encountered a problem and could not process your request"
	app.logger.Debug("internal error", "err", err.Error())
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *application) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Token")
	message := "invalid or missing authentication token"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}
