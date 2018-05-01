package todos

import (
	"errors"

	dbtodos "gotodo/database/todos"
)

var (
	// ErrDetailEmpty is returned when the detail param is empty.
	ErrDetailEmpty = errors.New("Detail parameter is empty")

	// ErrTodoNotFound is returned when a todo could not be found.
	ErrTodoNotFound = dbtodos.ErrTodoNotFound
)
