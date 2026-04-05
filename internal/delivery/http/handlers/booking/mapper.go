package booking

import (
	"time"

	"github.com/internships-backend/test-backend-the-new-day/internal/domain/entity"
)

func mapBookingToDto(booking *entity.Booking) BookingResponse {
	return BookingResponse{
		BookingID:      booking.BookingID.String(),
		SlotID:         booking.SlotID.String(),
		UserID:         booking.UserID.String(),
		Status:         string(booking.Status),
		ConferenceLink: booking.ConferenceLink,
		CreatedAt:      formatTimePtr(booking.CreatedAt),
	}
}

func mapBookingsToDto(bookings []*entity.Booking) []BookingResponse {
	res := make([]BookingResponse, len(bookings))
	for i, booking := range bookings {
		res[i] = mapBookingToDto(booking)
	}
	return res
}

func formatTimePtr(t time.Time) *string {
	val := t.UTC().Format(time.RFC3339)
	return &val
}
