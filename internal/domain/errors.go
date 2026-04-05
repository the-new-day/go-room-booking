package domain

import "errors"

var (
	ErrEmptyRoomName           = errors.New("room name cannot be empty")
	ErrNonPositiveRoomCapacity = errors.New("capacity must be positive")
)
