package data

import (
	"context"
	"database/sql"
	"errors"
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

func (m ChatModel) GetByUUID(uuid uuid.UUID) (*Chat, error) {
	chat := new(Chat)
	query := `
		SELECT uuid, created_at
		FROM chats
		WHERE uuid = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, uuid).Scan(&chat.UUID, &chat.CreatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errors.New("not found")
		default:
			return nil, err
		}
	}
	return chat, nil
}

type ChatMemberModel struct {
	DB *sql.DB
}

func (m ChatMemberModel) GetChatMembersByChat(chatUUID uuid.UUID) ([]*ChatMember, error) {
	query := `
		SELECT user_id, is_owner
		FROM chat_members
		WHERE chat_members.chat_uuid = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, chatUUID)
	if err != nil {
		return nil, err
	}

	chatMembers := []*ChatMember{}
	for rows.Next() {
		var member ChatMember
		err := rows.Scan(&member.UserID, &member.IsOwner)
		if err != nil {
			return nil, err
		}

		chatMembers = append(chatMembers, &member)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chatMembers, nil
}

type ChatMessageModel struct {
	DB *sql.DB
}
