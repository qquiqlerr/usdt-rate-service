package models

import "errors"

var (
	ErrIncompleteDepthData = errors.New("incomplete depth data")
	ErrInvalidTimestamp    = errors.New("invalid timestamp")
)
