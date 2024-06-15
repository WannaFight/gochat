package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type ChatMessage struct {
	ID     int64      `json:"id"`
	Text   string     `json:"text"`
	SentAt time.Time  `json:"sent_at"`
	Author ChatMember `json:"author"`
}

type ChatMessageModel struct {
	DB *sql.DB
}

func (m ChatMessageModel) Insert(message *ChatMessage) error {
	query := `
		INSERT INTO chat_messages (text, chat_member_id)
		VALUES ($1, $2)
		RETURNING id, sent_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, message.Text, message.Author.ID).Scan(
		&message.ID,
		&message.SentAt,
	)
}

func (m ChatMessageModel) GetAllByChat(chatUUID uuid.UUID) ([]*ChatMessage, error) {
	query := `
		SELECT msg.id, msg.text, msg.sent_at, mem.id, mem.is_owner, mem.user_id
		FROM chat_messages AS msg
		JOIN chat_members AS mem ON msg.chat_member_id = mem.id
		WHERE mem.chat_uuid = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, chatUUID)
	if err != nil {
		return nil, err
	}

	messages := []*ChatMessage{}

	for rows.Next() {
		var message ChatMessage
		err := rows.Scan(
			&message.ID,
			&message.Text,
			&message.SentAt,
			&message.Author.ID,
			&message.Author.IsOwner,
			&message.Author.UserID,
		)
		if err != nil {
			return nil, err
		}

		messages = append(messages, &message)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
