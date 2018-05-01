package services

import (
	"gotodo/database"
	"gotodo/services/members"
	"gotodo/services/todos"
)

// Services defines the services.
type Services struct {
	Members *members.Service
	Todos   *todos.Service
}

// New returns a new set of services.
func New(db *database.Database) *Services {
	return &Services{
		Members: members.New(db),
		Todos:   todos.New(db),
	}
}
