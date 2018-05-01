package todos

import "errors"

var (
	// ErrCreatedInvalid is returned when the created parameter is invalid.
	ErrCreatedInvalid = errors.New("Created parameter is invalid, must be a datetime string in RFC3339 format")

	// ErrLimitInvalid is returned when the limit parameter is invalid.
	ErrLimitInvalid = errors.New("Limit parameter is invalid, must be an integer")

	// ErrLimitMax is returned when the limit parameter is greater than the
	// maximum allowable limit.
	ErrLimitMax = errors.New("Limit parameter is greater than maximum allowable limit")
)
