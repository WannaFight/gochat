package main

import (
	"fmt"
	"net/http"
)

func (app *application) getChatMessages(w http.ResponseWriter, r *http.Request) {
	uuid, err := app.readUUIDParam(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err.Error())
		return
	}
	fmt.Fprintf(w, "listing all chat messages of chat with uuid=%s\n", uuid.String())
}

func (app *application) createChatMessage(w http.ResponseWriter, r *http.Request) {
	uuid, err := app.readUUIDParam(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err.Error())
		return
	}
	fmt.Fprintf(w, "add chat member to chat with uuid=%s\n", uuid.String())
}
