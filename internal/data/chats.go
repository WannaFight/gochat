package data

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	UUID      uuid.UUID `json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	Owner     User      `json:"owner"`
}

type ChatModel struct {
	DB *sql.DB
}
