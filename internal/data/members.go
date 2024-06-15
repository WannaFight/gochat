package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type ChatMember struct {
	UserID    int64     `json:"user_id"`
	IsOwner   bool      `json:"is_owner"`
	CreatedAt time.Time `json:"created_at"`
	Chat      *Chat     `json:"chat"`
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
