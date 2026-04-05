package room

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/api"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain/entity"
	"github.com/internships-backend/test-backend-the-new-day/pkg/logger/sl"
)

type CreateRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"`
	Capacity    *int    `json:"capacity" validate:"omitempty,min=1"`
}

type CreateResponse struct {
	RoomID      string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Capacity    *int    `json:"capacity"`
	CreatedAt   *string `json:"createdAt"`
}

type RoomCreator interface {
	CreateRoom(ctx context.Context, name string, description *string, capacity *int) (*entity.Room, error)
}

func NewCreateHandler(logger *slog.Logger, roomCreator RoomCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, ok := api.DecodeRequest[CreateRequest](logger, w, r)
		if !ok {
			return
		}

		room, err := roomCreator.CreateRoom(r.Context(), req.Name, req.Description, req.Capacity)
		if err != nil {
			logger.Error("failed to create room", sl.Err(err))

			api.SendInternalError(w, r, "failed to register user")
			return
		}

		api.SendCreated(w, r, CreateResponse{
			RoomID:      room.RoomID.String(),
			Name:        room.Name,
			Description: room.Description,
			Capacity:    room.Capacity,
			CreatedAt:   nil,
		})
	}
}
