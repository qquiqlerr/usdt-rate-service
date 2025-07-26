package models

import "errors"

var (
	// ErrIncompleteDepthData indicates that the depth data (asks or bids) is incomplete.
	ErrIncompleteDepthData = errors.New("incomplete depth data")
	// ErrInvalidTimestamp indicates that the timestamp is invalid.
	ErrInvalidTimestamp = errors.New("invalid timestamp")
)
