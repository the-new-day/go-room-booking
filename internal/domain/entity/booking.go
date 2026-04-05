package entity

import (
	"time"

	"github.com/google/uuid"
)

type BookingStatus string

const (
	BookingActive    BookingStatus = "active"
	BookingCancelled BookingStatus = "cancelled"
)

type Booking struct {
	BookingID      uuid.UUID     `db:"booking_id"`
	SlotID         uuid.UUID     `db:"slot_id"`
	UserID         uuid.UUID     `db:"user_id"`
	Status         BookingStatus `db:"status"`
	ConferenceLink *string       `db:"conference_link"`
	CreatedAt      time.Time     `db:"created_at"`
}
