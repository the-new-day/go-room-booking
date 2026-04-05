package slot

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/api"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain/entity"
	"github.com/internships-backend/test-backend-the-new-day/pkg/logger/sl"
)

type slotDto struct {
	SlotID string `json:"id"`
	RoomID string `json:"roomId"`
	Start  string `json:"start"`
	End    string `json:"end"`
}

type ListResponse struct {
	Slots []slotDto `json:"slots"`
}

type SlotLister interface {
	ListAvailableSlots(ctx context.Context, roomID string, date time.Time) ([]*entity.Slot, error)
}

func NewListHandler(logger *slog.Logger, slotLister SlotLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomID := chi.URLParam(r, "roomId")
		if _, err := uuid.Parse(roomID); err != nil {
			api.SendBadRequest(w, r, "invalid room id")
			return
		}

		dateStr := r.URL.Query().Get("date")
		if dateStr == "" {
			api.SendBadRequest(w, r, "date is required")
			return
		}

		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			api.SendBadRequest(w, r, "invalid date")
			return
		}

		slots, err := slotLister.ListAvailableSlots(r.Context(), roomID, date)
		if errors.Is(err, domain.ErrRoomNotFound) {
			api.SendRoomNotFound(w, r, "room not found")
			return
		} else if err != nil {
			logger.Error("failed to list slots", sl.Err(err))
			api.SendInternalError(w, r, "failed to list slots")
			return
		}

		api.SendOK(w, r, ListResponse{
			Slots: mapSlotsToDto(slots),
		})
	}
}

func mapSlotsToDto(slots []*entity.Slot) []slotDto {
	res := make([]slotDto, len(slots))
	for i, slot := range slots {
		res[i] = slotDto{
			SlotID: slot.SlotID.String(),
			RoomID: slot.RoomID.String(),
			Start:  slot.StartAt.UTC().Format(time.RFC3339),
			End:    slot.EndAt.UTC().Format(time.RFC3339),
		}
	}
	return res
}
