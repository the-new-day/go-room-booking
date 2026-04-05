package room

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/api"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain/entity"
	"github.com/internships-backend/test-backend-the-new-day/pkg/logger/sl"
)

type roomDto struct {
	RoomID      string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Capacity    *int    `json:"capacity"`
	CreatedAt   *string `json:"createdAt"`
}

type ListResponse struct {
	Rooms []roomDto `json:"rooms"`
}

type RoomLister interface {
	ListRooms(ctx context.Context) ([]*entity.Room, error)
}

func NewListHandler(logger *slog.Logger, roomLister RoomLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rooms, err := roomLister.ListRooms(r.Context())

		if err != nil {
			logger.Error("failed to list rooms", sl.Err(err))

			api.SendInternalError(w, r, "failed to list rooms")
			return
		}

		api.SendOK(w, r, ListResponse{
			Rooms: mapRoomsToDto(rooms),
		})
	}
}

func mapRoomsToDto(rooms []*entity.Room) []roomDto {
	res := make([]roomDto, len(rooms))
	for i, room := range rooms {
		res[i].RoomID = room.RoomID.String()
		res[i].Name = room.Name
		res[i].Description = room.Description
		res[i].Capacity = room.Capacity
		res[i].CreatedAt = nil
	}
	return res
}
