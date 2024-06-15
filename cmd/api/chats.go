package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/WannaFight/gochat/internal/data"
)

func (app *application) listChatHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	chats, err := app.models.Chats.GetAllByUser(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"chats": chats}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createChatHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	chat := &data.Chat{
		Name: input.Name,
	}

	err = app.models.Chats.Insert(chat)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Set current user as chat owner.
	user := app.contextGetUser(r)
	chatOwner := &data.ChatMember{
		UserID:  user.ID,
		IsOwner: true,
		Chat:    chat,
	}
	err = app.models.ChatMembers.Insert(chatOwner)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"chat": chat}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getChatHandler(w http.ResponseWriter, r *http.Request) {
	uuid, err := app.readUUIDParam(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err.Error())
		return
	}

	chat, err := app.models.Chats.GetByUUID(uuid)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"chat": chat}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getChatMembers(w http.ResponseWriter, r *http.Request) {
	uuid, err := app.readUUIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	members, err := app.models.ChatMembers.GetAllByChat(uuid)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"chat_members": members}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createChatMember(w http.ResponseWriter, r *http.Request) {
	uuid, err := app.readUUIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	var input struct {
		Username string `json:"username"`
		IsOwner  bool   `json:"is_owner"`
	}

	// Getting chat
	chat, err := app.models.Chats.GetByUUID(uuid)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.models.Users.GetByUsername(input.Username)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	member := &data.ChatMember{
		UserID:  user.ID,
		IsOwner: input.IsOwner,
		Chat:    chat,
	}

	err = app.models.ChatMembers.Insert(member)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"chat_member": member}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
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

func (app *application) createChatMessage(w http.ResponseWriter, r *http.Request) {
	uuid, err := app.readUUIDParam(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err.Error())
		return
	}
	fmt.Fprintf(w, "add chat member to chat with uuid=%s\n", uuid.String())
}
