package domain

import "errors"

var (
	ErrEmptyRoomName           = errors.New("room name cannot be empty")
	ErrNonPositiveRoomCapacity = errors.New("capacity must be positive")

	ErrUserWithEmailAlreadyExists = errors.New("user with provided email already exists")

	ErrEmailNotFound   = errors.New("user with provided email not found")
	ErrInvalidPassword = errors.New("password is invalid")

	ErrInvalidRole = errors.New("invalid role")

	ErrRoomNotFound          = errors.New("room not found")
	ErrScheduleAlreadyExists = errors.New("schedule already exists")
	ErrInvalidDaysOfWeek     = errors.New("days of week must be between 1 and 7")
	ErrInvalidTimeRange      = errors.New("start time must be before end time")

	ErrSlotNotFound      = errors.New("slot not found")
	ErrSlotAlreadyBooked = errors.New("slot already booked")
	ErrSlotInPast        = errors.New("slot is in the past")

	ErrBookingNotFound = errors.New("booking not found")
	ErrForbidden       = errors.New("forbidden")
)
