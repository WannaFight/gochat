package main

import (
	"fmt"
	"net/http"
)

func (app *application) listChatHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "listing all user chats")
}

func (app *application) createChatHandler(w http.ResponseWriter, r *http.Request) {
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

func (app *application) getChatMessages(w http.ResponseWriter, r *http.Request) {
	uuid, err := app.readUUIDParam(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err.Error())
		return
	}
	fmt.Fprintf(w, "listing all chat messages of chat with uuid=%s\n", uuid.String())
}
