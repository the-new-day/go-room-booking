package entity

import (
	"time"

	"github.com/google/uuid"
)

type Slot struct {
	SlotID  uuid.UUID `db:"slot_id"`
	RoomID  uuid.UUID `db:"room_id"`
	StartAt time.Time `db:"start_at"`
	EndAt   time.Time `db:"end_at"`
}
