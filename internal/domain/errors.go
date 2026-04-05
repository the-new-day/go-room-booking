package domain

import "errors"

var (
	ErrEmptyRoomName           = errors.New("room name cannot be empty")
	ErrNonPositiveRoomCapacity = errors.New("capacity must be positive")

	ErrUserWithEmailAlreadyExists = errors.New("user with provided email already exists")

	ErrEmailNotFound   = errors.New("user with provided email not found")
	ErrInvalidPassword = errors.New("password is invalid")

	ErrInvalidRole = errors.New("invalid role")
)
