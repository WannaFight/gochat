package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Users        UserModel
	Chats        ChatModel
	ChatMembers  ChatMemberModel
	ChatMessages ChatMessageModel
	Tokens       TokenModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:        UserModel{DB: db},
		Chats:        ChatModel{DB: db},
		ChatMembers:  ChatMemberModel{DB: db},
		ChatMessages: ChatMessageModel{DB: db},
		Tokens:       TokenModel{DB: db},
	}
}
