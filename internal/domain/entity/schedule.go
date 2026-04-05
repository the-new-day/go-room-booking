package entity

import (
	"time"

	"github.com/google/uuid"
)

const SlotDuration = 30 * time.Minute

type Schedule struct {
	ScheduleID uuid.UUID    `db:"schedule_id"`
	RoomID     uuid.UUID    `db:"room_id"`
	Weekday    time.Weekday `db:"weekday"`
	StartAt    time.Time    `db:"start_at"`
}
