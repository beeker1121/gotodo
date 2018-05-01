package todos

import "errors"

var (
	// ErrTodoNotFound is returned when a todo could not be found.
	ErrTodoNotFound = errors.New("Todo could not be found")
)
