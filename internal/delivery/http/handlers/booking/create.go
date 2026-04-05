package booking

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/api"
	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/middleware"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain/entity"
	"github.com/internships-backend/test-backend-the-new-day/pkg/logger/sl"
)

type CreateRequest struct {
	SlotID               string `json:"slotId" validate:"required"`
	CreateConferenceLink bool   `json:"createConferenceLink"`
}

type BookingResponse struct {
	BookingID      string  `json:"id"`
	SlotID         string  `json:"slotId"`
	UserID         string  `json:"userId"`
	Status         string  `json:"status"`
	ConferenceLink *string `json:"conferenceLink"`
	CreatedAt      *string `json:"createdAt"`
}

type CreateResponse struct {
	Booking BookingResponse `json:"booking"`
}

type BookingCreator interface {
	CreateBooking(ctx context.Context, slotID, userID string, createConferenceLink bool) (*entity.Booking, error)
}

func NewCreateHandler(logger *slog.Logger, bookingCreator BookingCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, ok := api.DecodeRequest[CreateRequest](logger, w, r)
		if !ok {
			return
		}

		if _, err := uuid.Parse(req.SlotID); err != nil {
			api.SendBadRequest(w, r, "invalid slot id")
			return
		}

		userID, ok := r.Context().Value(middleware.UserIDKey).(string)
		if !ok || userID == "" {
			api.SendUnauthorized(w, r, "unauthorized")
			return
		}

		booking, err := bookingCreator.CreateBooking(r.Context(), req.SlotID, userID, req.CreateConferenceLink)
		if errors.Is(err, domain.ErrSlotNotFound) {
			api.SendSlotNotFound(w, r, "slot not found")
			return
		} else if errors.Is(err, domain.ErrSlotAlreadyBooked) {
			api.SendSlotAlreadyBooked(w, r, "slot is already booked")
			return
		} else if errors.Is(err, domain.ErrSlotInPast) {
			api.SendBadRequest(w, r, "slot start time is in the past")
			return
		} else if err != nil {
			logger.Error("failed to create booking", sl.Err(err))
			api.SendInternalError(w, r, "failed to create booking")
			return
		}

		api.SendCreated(w, r, CreateResponse{
			Booking: mapBookingToDto(booking),
		})
	}
}
