package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID     uuid.UUID `json:"id"`
	Text   string    `json:"text"`
	SentAt time.Time `json:"sent_at"`
	Author User      `json:"author"`
}

type MessageModel struct {
	DB *sql.DB
}

func (m MessageModel) Insert(message *Message) error {
	query := `
		INSERT INTO messages (text, user_id)
		VALUES ($1, $2)
		RETURNING id, sent_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, message.Text, message.Author.ID).Scan(
		&message.ID,
		&message.SentAt,
	)
}

func (m MessageModel) GetAll() ([]*Message, error) {
	query := `
		SELECT m.id, m.text, m.sent_at, u.username
		FROM messages AS m
		JOIN users AS u ON u.id = m.user_id
		ORDER BY m.created_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	messages := []*Message{}

	for rows.Next() {
		var message Message
		err := rows.Scan(
			&message.ID,
			&message.Text,
			&message.SentAt,
			&message.Author.Username,
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
