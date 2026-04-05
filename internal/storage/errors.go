package storage

import "errors"

var (
	ErrNotFound              = errors.New("not found")
	ErrAlreadyExists         = errors.New("already exists")
	ErrDuplicateEmail        = errors.New("user with this email already exists")
	ErrScheduleAlreadyExists = errors.New("schedule already exists for this room")
	ErrSlotOverlap           = errors.New("slot overlaps with existing slot")
	ErrSlotAlreadyBooked     = errors.New("slot is already booked")
	ErrInvalidTimeRange      = errors.New("start time must be before end time")
	ErrInvalidDaysOfWeek     = errors.New("days of week must be between 1 and 7")
	ErrInvalidRole           = errors.New("role must be admin or user")
	ErrInvalidCapacity       = errors.New("capacity must be positive")
)
