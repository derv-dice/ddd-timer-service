package http

import "errors"

const messageJsonKey = "message"

var (
	ErrInvalidDateFormat = errors.New("invalid date format")
	ErrBadDates          = errors.New("'from' must be before 'to'")
)
