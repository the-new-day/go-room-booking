package entity

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	RoomID      uuid.UUID `db:"room_id"`
	Name        string    `db:"name"`
	Description *string   `db:"description"`
	Capacity    *int      `db:"capacity"`
	CreatedAt   time.Time `db:"created_at"`
}
