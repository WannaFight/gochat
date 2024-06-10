package data

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	UUID        uuid.UUID    `json:"uuid"`
	CreatedAt   time.Time    `json:"created_at"`
	ChatMembers []ChatMember `json:"chat_members"`
}

type ChatMember struct {
	UserID  int64 `json:"user_id"`
	IsOwner bool  `json:"is_owner"`
}

type ChatMessage struct {
	Text   string     `json:"text"`
	SentAt time.Time  `json:"sent_at"`
	Author ChatMember `json:"author"`
}

type ChatModel struct {
	DB *sql.DB
}

type ChatMemberModel struct {
	DB *sql.DB
}

type ChatMessageModel struct {
	DB *sql.DB
}
