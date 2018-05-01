package members

import "errors"

var (
	// ErrMemberNotFound is returned when a member could not be found.
	ErrMemberNotFound = errors.New("Member could not be found")
)
