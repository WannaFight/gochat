package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type ChatMember struct {
	ID        int64     `json:"id"`
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
		INSERT INTO chat_members (chat_id, user_id, is_owner)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`
	args := []any{chatMember.Chat.UUID, chatMember.UserID, chatMember.IsOwner}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&chatMember.ID, &chatMember.CreatedAt)
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

func (m ChatMemberModel) GetByIDAndChat(userID int64, chatUUID uuid.UUID) (*ChatMember, error) {
	query := `
		SELECT id, user_id, is_owner
		FROM chat_members
		WHERE user_id = $1 AND chat_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	member := new(ChatMember)
	err := m.DB.QueryRowContext(ctx, query, userID, chatUUID).Scan(
		&member.ID,
		&member.UserID,
		&member.IsOwner,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return member, nil
}

func (m ChatMemberModel) GetAllByChat(chatUUID uuid.UUID) ([]*ChatMember, error) {
	query := `
		SELECT id, user_id, is_owner
		FROM chat_members
		WHERE chat_id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, chatUUID)
	if err != nil {
		return nil, err
	}

	chatMembers := []*ChatMember{}
	for rows.Next() {
		var member ChatMember
		err := rows.Scan(&member.ID, &member.UserID, &member.IsOwner)
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
