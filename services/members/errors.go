package members

import (
	"errors"

	dbmembers "gotodo/database/members"
)

var (
	// ErrEmailEmpty is returned when the email param is empty.
	ErrEmailEmpty = errors.New("Email parameter is empty")

	// ErrEmailExists is returned when the email already exists.
	ErrEmailExists = errors.New("Email already exists")

	// ErrPassword is returned when the password is in an invalid format.
	ErrPassword = errors.New("Password must be at least 8 characters")

	// ErrInvalidLogin is returned when the email and/or password used
	// with login is invalid.
	ErrInvalidLogin = errors.New("Email and/or password is invalid")

	// ErrMemberNotFound is returned when a member could not be found.
	ErrMemberNotFound = dbmembers.ErrMemberNotFound
)
