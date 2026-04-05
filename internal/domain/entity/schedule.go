package entity

import (
	"time"

	"github.com/google/uuid"
)

const SlotDuration = 30 * time.Minute

type Weekday int

const (
	Monday Weekday = iota + 1
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

type Schedule struct {
	ScheduleID uuid.UUID `db:"schedule_id"`
	RoomID     uuid.UUID `db:"room_id"`
	Weekdays   []Weekday `db:"weekdays"`
	StartAt    time.Time `db:"start_at"`
	EndAt      time.Time `db:"end_at"`
	CreatedAt  time.Time `db:"created_at"`
}
