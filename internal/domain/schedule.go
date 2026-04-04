package domain

import (
	"time"

	"github.com/google/uuid"
)

type Schedule struct {
	ScheduleID uuid.UUID    `db:"schedule_id"`
	RoomID     uuid.UUID    `db:"room_id"`
	Weekday    time.Weekday `db:"weekday"`
	StartAt    string       `db:"start_at"`
	EndAt      string       `db:"end_at"`
}
