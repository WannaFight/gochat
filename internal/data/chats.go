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

type ChatMember struct {
	UserID    int64     `json:"user_id"`
	IsOwner   bool      `json:"is_owner"`
	CreatedAt time.Time `json:"created_at"`
	Chat      *Chat     `json:"chat"`
}

type ChatMessage struct {
	Text   string     `json:"text"`
	SentAt time.Time  `json:"sent_at"`
	Author ChatMember `json:"author"`
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
		SELECT chats.uuid, chats.name, chats.created_at
		FROM chats
		JOIN chat_members ON chats.uuid = chat_members.chat_uuid
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

type ChatMemberModel struct {
	DB *sql.DB
}

func (m ChatMemberModel) Insert(chatMember *ChatMember) error {
	query := `
		INSERT INTO chat_members (chat_uuid, user_id, is_owner)
		VALUES ($1, $2, $3)
		RETURNING created_at`
	args := []any{chatMember.Chat.UUID, chatMember.UserID, chatMember.IsOwner}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&chatMember.CreatedAt)
	if err != nil {
		switch err.Error() {
		// User with userID not found.
		case `pq: insert or update on table "chat_members" violates foreign key constraint "chat_members_user_id_fkey"`:
			return ErrRecordNotFound
		default:
			return err
		}
	}
	return nil
}

func (m ChatMemberModel) GetAllByChat(chatUUID uuid.UUID) ([]*ChatMember, error) {
	query := `
		SELECT user_id, is_owner
		FROM chat_members
		WHERE chat_uuid = $1`
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
