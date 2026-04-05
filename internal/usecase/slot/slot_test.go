package slot

import (
	"errors"
	"testing"
	"time"

	"github.com/internships-backend/test-backend-the-new-day/internal/domain"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain/entity"
	"github.com/internships-backend/test-backend-the-new-day/internal/usecase/slot/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSlotUseCase_ListAvailableSlots(t *testing.T) {
	errDummy := errors.New("db error")
	date := time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		roomID   string
		setup    func(roomRepo *mocks.MockRoomRepository, scheduleRepo *mocks.MockScheduleRepository, slotRepo *mocks.MockSlotRepository)
		wantErr  error
		wantSize int
	}{
		{
			name:   "room not found",
			roomID: "room-1",
			setup: func(roomRepo *mocks.MockRoomRepository, scheduleRepo *mocks.MockScheduleRepository, slotRepo *mocks.MockSlotRepository) {
				roomRepo.EXPECT().Exists(mock.Anything, "room-1").Return(false, nil)
			},
			wantErr: domain.ErrRoomNotFound,
		},
		{
			name:   "no schedule",
			roomID: "room-1",
			setup: func(roomRepo *mocks.MockRoomRepository, scheduleRepo *mocks.MockScheduleRepository, slotRepo *mocks.MockSlotRepository) {
				roomRepo.EXPECT().Exists(mock.Anything, "room-1").Return(true, nil)
				scheduleRepo.EXPECT().ExistsByRoomID(mock.Anything, "room-1").Return(false, nil)
			},
			wantSize: 0,
		},
		{
			name:   "success",
			roomID: "room-1",
			setup: func(roomRepo *mocks.MockRoomRepository, scheduleRepo *mocks.MockScheduleRepository, slotRepo *mocks.MockSlotRepository) {
				roomRepo.EXPECT().Exists(mock.Anything, "room-1").Return(true, nil)
				scheduleRepo.EXPECT().ExistsByRoomID(mock.Anything, "room-1").Return(true, nil)
				slotRepo.EXPECT().ListAvailableByRoomAndDate(mock.Anything, "room-1", date).
					Return([]*entity.Slot{{}, {}}, nil)
			},
			wantSize: 2,
		},
		{
			name:   "schedule repo error",
			roomID: "room-1",
			setup: func(roomRepo *mocks.MockRoomRepository, scheduleRepo *mocks.MockScheduleRepository, slotRepo *mocks.MockSlotRepository) {
				roomRepo.EXPECT().Exists(mock.Anything, "room-1").Return(true, nil)
				scheduleRepo.EXPECT().ExistsByRoomID(mock.Anything, "room-1").Return(false, errDummy)
			},
			wantErr: errDummy,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			roomRepo := mocks.NewMockRoomRepository(t)
			scheduleRepo := mocks.NewMockScheduleRepository(t)
			slotRepo := mocks.NewMockSlotRepository(t)

			if tt.setup != nil {
				tt.setup(roomRepo, scheduleRepo, slotRepo)
			}

			uc := New(roomRepo, scheduleRepo, slotRepo)

			slots, err := uc.ListAvailableSlots(t.Context(), tt.roomID, date)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, slots, tt.wantSize)
		})
	}
}
