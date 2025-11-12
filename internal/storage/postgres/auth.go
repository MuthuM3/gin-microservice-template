package postgres

import "database/sql"

type AuthStore struct {
	db *sql.DB
}

func NewAuthStore(db *sql.DB) *AuthStore {
	return &AuthStore{
		db: db,
	}
}

