package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Users    UserModel
	Messages MessageModel
	Tokens   TokenModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:    UserModel{DB: db},
		Messages: MessageModel{DB: db},
		Tokens:   TokenModel{DB: db},
	}
}
