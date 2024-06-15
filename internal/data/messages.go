package data

import (
	"database/sql"
	"time"
)

type ChatMessage struct {
	Text   string     `json:"text"`
	SentAt time.Time  `json:"sent_at"`
	Author ChatMember `json:"author"`
}

type ChatMessageModel struct {
	DB *sql.DB
}
