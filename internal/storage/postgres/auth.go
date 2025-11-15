package postgres

import "database/sql"

type AuthStore struct {
	db    *sql.DB
	store *Store
}

func NewAuthStore(db *sql.DB, store *Store) *AuthStore {
	return &AuthStore{
		db:    db,
		store: store,
	}
}
