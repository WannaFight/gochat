package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	UUID      uuid.UUID `json:"uuid"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type ChatModel struct {
	DB *sql.DB
}

func (m ChatModel) Insert(chat *Chat) error {
	query := `
		INSERT INTO chats (name)
		VALUES ($1)
		RETURNING uuid, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, chat.Name).Scan(
		&chat.UUID,
		&chat.CreatedAt,
	)
}

func (m ChatModel) GetAllByUser(userID int64) ([]*Chat, error) {
	query := `
		SELECT chats.id, chats.name, chats.created_at
		FROM chats
		JOIN chat_members ON chats.id = chat_members.chat_id
		WHERE chat_members.user_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	chats := []*Chat{}
	for rows.Next() {
		var chat Chat
		err := rows.Scan(&chat.UUID, &chat.Name, &chat.CreatedAt)
		if err != nil {
			return nil, err
		}

		chats = append(chats, &chat)

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chats, nil
}

func (m ChatModel) GetByUUID(uuid uuid.UUID) (*Chat, error) {
	chat := new(Chat)
	query := `
		SELECT uuid, name, created_at
		FROM chats
		WHERE uuid = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, uuid).Scan(
		&chat.UUID,
		&chat.Name,
		&chat.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return chat, nil
}
