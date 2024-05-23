package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (app *application) readUUIDParam(r *http.Request) (uuid.UUID, error) {
	return uuid.Parse(r.PathValue("uuid"))
}
