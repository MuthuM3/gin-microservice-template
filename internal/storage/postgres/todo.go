package postgres

import "database/sql"

type TodoStore struct {
	db    *sql.DB
	store *Store
}

func newTodoStore(db *sql.DB, store *Store) *TodoStore {
	return &TodoStore{
		db:    db,
		store: store,
	}
}
