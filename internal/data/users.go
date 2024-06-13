package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateUsername = errors.New("duplicate username")
)

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username"`
	Password  password  `json:"-"`
}

type password struct {
	plain *string
	hash  []byte
}

func (p *password) Set(plainPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), 12)

	if err != nil {
		return err
	}

	p.plain = &plainPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plainPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainPassword))

	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, err
}

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) Insert(user *User) error {
	query := `
		INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING id, created_at`
	args := []any{user.Username, user.Password.hash}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}

func (m *UserModel) GetByUsername(username string) (*User, error) {
	user := new(User)
	query := `
		SELECT id, username, created_at
		FROM users
		WHERE username = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}
