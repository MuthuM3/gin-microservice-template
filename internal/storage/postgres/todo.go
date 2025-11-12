package postgres

import "database/sql"

type TodoStore struct {
	db *sql.DB
}

func newTodoStore(db *sql.DB) *TodoStore {
	return &TodoStore{
		db: db,
	}
}