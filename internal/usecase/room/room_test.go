package room

import (
	"strings"
	"testing"

	"github.com/internships-backend/test-backend-the-new-day/internal/domain"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain/entity"
	"github.com/internships-backend/test-backend-the-new-day/internal/usecase/room/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ptr[T any](v T) *T {
	return &v
}

func TestRoomUseCase_CreateRoom(t *testing.T) {
	type room struct {
		name        string
		description *string
		capacity    *int
	}

	tests := []struct {
		name         string
		room         room
		wantErr      error
		wantRoomName string
	}{
		{
			name: "success",
			room: room{
				name:        "room 1",
				description: ptr("123"),
				capacity:    ptr(5),
			},
			wantRoomName: "room 1",
		},
		{
			name: "negative capacity",
			room: room{
				name:        "room 1",
				description: ptr("123"),
				capacity:    ptr(-5),
			},
			wantErr: domain.ErrNonPositiveRoomCapacity,
		},
		{
			name: "zero capacity",
			room: room{
				name:        "room 1",
				description: ptr("123"),
				capacity:    ptr(0),
			},
			wantErr: domain.ErrNonPositiveRoomCapacity,
		},
		{
			name: "success with empty description",
			room: room{
				name:        "room 1",
				description: ptr(""),
				capacity:    ptr(5),
			},
			wantRoomName: "room 1",
		},
		{
			name: "success with null description",
			room: room{
				name:     "room 1",
				capacity: ptr(5),
			},
			wantRoomName: "room 1",
		},
		{
			name: "success with null capacity",
			room: room{
				name: "room 1",
			},
			wantRoomName: "room 1",
		},
		{
			name: "empty room name",
			room: room{
				name: "",
			},
			wantErr: domain.ErrEmptyRoomName,
		},
		{
			name: "empty trimmed room name",
			room: room{
				name: "   ",
			},
			wantErr: domain.ErrEmptyRoomName,
		},
		{
			name: "trim spaces in room name",
			room: room{
				name: "    room 1     ",
			},
			wantRoomName: "room 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := mocks.NewMockRoomRepository(t)
			returnRoom := &entity.Room{
				Name:        strings.TrimSpace(tt.room.name),
				Description: tt.room.description,
				Capacity:    tt.room.capacity,
			}

			mockRepo.EXPECT().
				Create(t.Context(), strings.TrimSpace(tt.room.name), tt.room.description, tt.room.capacity).
				Return(returnRoom, tt.wantErr).
				Maybe()

			uc := New(mockRepo)

			room, err := uc.CreateRoom(t.Context(), tt.room.name, tt.room.description, tt.room.capacity)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantRoomName, room.Name)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
