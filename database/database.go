package database

import (
	"database/sql"

	"gotodo/database/members"
	"gotodo/database/todos"
)

// Database defines the database.
type Database struct {
	Members *members.Database
	Todos   *todos.Database
}

// New returns a new database.
func New(db *sql.DB) *Database {
	return &Database{
		Members: members.New(db),
		Todos:   todos.New(db),
	}
}
