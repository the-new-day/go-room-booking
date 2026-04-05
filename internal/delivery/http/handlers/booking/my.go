package booking

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/api"
	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/middleware"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain/entity"
	"github.com/internships-backend/test-backend-the-new-day/pkg/logger/sl"
)

type MyResponse struct {
	Bookings []BookingResponse `json:"bookings"`
}

type MyBookingLister interface {
	ListMyBookings(ctx context.Context, userID string) ([]*entity.Booking, error)
}

func NewMyHandler(logger *slog.Logger, bookingLister MyBookingLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middleware.UserIDKey).(string)
		if !ok || userID == "" {
			api.SendUnauthorized(w, r, "unauthorized")
			return
		}

		bookings, err := bookingLister.ListMyBookings(r.Context(), userID)
		if err != nil {
			logger.Error("failed to list user bookings", sl.Err(err))
			api.SendInternalError(w, r, "failed to list bookings")
			return
		}

		api.SendOK(w, r, MyResponse{
			Bookings: mapBookingsToDto(bookings),
		})
	}
}
