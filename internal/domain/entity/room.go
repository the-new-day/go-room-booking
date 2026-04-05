package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrEmptyRoomName           = errors.New("room name cannot be empty")
	ErrNonPositiveRoomCapacity = errors.New("capacity must be positive")
)

type Room struct {
	RoomID      uuid.UUID `db:"room_id"`
	Name        string    `db:"name"`
	Description *string   `db:"description"`
	Capacity    *int      `db:"capacity"`
	CreatedAt   time.Time `db:"created_at"`
}

func (r *Room) Validate() error {
	if r.Name == "" {
		return ErrEmptyRoomName
	}
	if r.Capacity != nil && *r.Capacity <= 0 {
		return ErrNonPositiveRoomCapacity
	}
	return nil
}
