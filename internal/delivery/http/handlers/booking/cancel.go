package booking

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/api"
	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/middleware"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain/entity"
	"github.com/internships-backend/test-backend-the-new-day/pkg/logger/sl"
)

type CancelResponse struct {
	Booking BookingResponse `json:"booking"`
}

type BookingCanceler interface {
	CancelBooking(ctx context.Context, bookingID, userID string) (*entity.Booking, error)
}

func NewCancelHandler(logger *slog.Logger, bookingCanceler BookingCanceler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bookingID := chi.URLParam(r, "bookingId")
		if _, err := uuid.Parse(bookingID); err != nil {
			api.SendBadRequest(w, r, "invalid booking id")
			return
		}

		userID, ok := r.Context().Value(middleware.UserIDKey).(string)
		if !ok || userID == "" {
			api.SendUnauthorized(w, r, "unauthorized")
			return
		}

		booking, err := bookingCanceler.CancelBooking(r.Context(), bookingID, userID)
		if errors.Is(err, domain.ErrBookingNotFound) {
			api.SendBookingNotFound(w, r, "booking not found")
			return
		} else if errors.Is(err, domain.ErrForbidden) {
			api.SendForbidden(w, r, "cannot cancel another user's booking")
			return
		} else if err != nil {
			logger.Error("failed to cancel booking", sl.Err(err))
			api.SendInternalError(w, r, "failed to cancel booking")
			return
		}

		api.SendOK(w, r, CancelResponse{
			Booking: mapBookingToDto(booking),
		})
	}
}
