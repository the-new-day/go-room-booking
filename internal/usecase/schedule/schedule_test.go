package schedule

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain/entity"
	"github.com/internships-backend/test-backend-the-new-day/internal/storage"
	"github.com/internships-backend/test-backend-the-new-day/internal/usecase/schedule/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestScheduleUseCase_CreateSchedule(t *testing.T) {
	errDummy := errors.New("db error")

	tests := []struct {
		name      string
		roomID    string
		weekdays  []int
		startTime string
		endTime   string
		setup     func(roomRepo *mocks.MockRoomRepository, scheduleRepo *mocks.MockScheduleRepository, slotRepo *mocks.MockSlotRepository)
		wantErr   error
		wantCalls bool
	}{
		{
			name:      "room not found",
			roomID:    "room-1",
			weekdays:  []int{1},
			startTime: "09:00",
			endTime:   "10:00",
			setup: func(roomRepo *mocks.MockRoomRepository, scheduleRepo *mocks.MockScheduleRepository, slotRepo *mocks.MockSlotRepository) {
				roomRepo.EXPECT().Exists(mock.Anything, "room-1").Return(false, nil)
			},
			wantErr: domain.ErrRoomNotFound,
		},
		{
			name:      "empty weekdays",
			roomID:    "room-1",
			weekdays:  []int{},
			startTime: "09:00",
			endTime:   "10:00",
			setup: func(roomRepo *mocks.MockRoomRepository, scheduleRepo *mocks.MockScheduleRepository, slotRepo *mocks.MockSlotRepository) {
				roomRepo.EXPECT().Exists(mock.Anything, "room-1").Return(true, nil)
			},
			wantErr: domain.ErrInvalidDaysOfWeek,
		},
		{
			name:      "invalid weekday",
			roomID:    "room-1",
			weekdays:  []int{0},
			startTime: "09:00",
			endTime:   "10:00",
			setup: func(roomRepo *mocks.MockRoomRepository, scheduleRepo *mocks.MockScheduleRepository, slotRepo *mocks.MockSlotRepository) {
				roomRepo.EXPECT().Exists(mock.Anything, "room-1").Return(true, nil)
			},
			wantErr: domain.ErrInvalidDaysOfWeek,
		},
		{
			name:      "invalid time range",
			roomID:    "room-1",
			weekdays:  []int{1},
			startTime: "10:00",
			endTime:   "09:00",
			setup: func(roomRepo *mocks.MockRoomRepository, scheduleRepo *mocks.MockScheduleRepository, slotRepo *mocks.MockSlotRepository) {
				roomRepo.EXPECT().Exists(mock.Anything, "room-1").Return(true, nil)
			},
			wantErr: domain.ErrInvalidTimeRange,
		},
		{
			name:      "schedule already exists",
			roomID:    "room-1",
			weekdays:  []int{1},
			startTime: "09:00",
			endTime:   "10:00",
			setup: func(roomRepo *mocks.MockRoomRepository, scheduleRepo *mocks.MockScheduleRepository, slotRepo *mocks.MockSlotRepository) {
				roomRepo.EXPECT().Exists(mock.Anything, "room-1").Return(true, nil)
				scheduleRepo.EXPECT().Create(mock.Anything, "room-1", []int{1}, "09:00", "10:00").
					Return(nil, storage.ErrScheduleAlreadyExists)
			},
			wantErr: domain.ErrScheduleAlreadyExists,
		},
		{
			name:      "success",
			roomID:    "room-1",
			weekdays:  []int{1, 2, 3, 4, 5, 6, 7},
			startTime: "09:00",
			endTime:   "11:00",
			setup: func(roomRepo *mocks.MockRoomRepository, scheduleRepo *mocks.MockScheduleRepository, slotRepo *mocks.MockSlotRepository) {
				roomRepo.EXPECT().Exists(mock.Anything, "room-1").Return(true, nil)
				scheduleRepo.EXPECT().Create(mock.Anything, "room-1", []int{1, 2, 3, 4, 5, 6, 7}, "09:00", "11:00").
					Return(&entity.Schedule{ScheduleID: testUUID(), RoomID: testUUID()}, nil)
				slotRepo.EXPECT().
					Create(mock.Anything, "room-1", mock.Anything, mock.Anything).
					Return(nil).
					Times(32)
			},
		},
		{
			name:      "storage invalid days of week",
			roomID:    "room-1",
			weekdays:  []int{1},
			startTime: "09:00",
			endTime:   "10:00",
			setup: func(roomRepo *mocks.MockRoomRepository, scheduleRepo *mocks.MockScheduleRepository, slotRepo *mocks.MockSlotRepository) {
				roomRepo.EXPECT().Exists(mock.Anything, "room-1").Return(true, nil)
				scheduleRepo.EXPECT().Create(mock.Anything, "room-1", []int{1}, "09:00", "10:00").
					Return(nil, storage.ErrInvalidDaysOfWeek)
			},
			wantErr: domain.ErrInvalidDaysOfWeek,
		},
		{
			name:      "storage invalid time range",
			roomID:    "room-1",
			weekdays:  []int{1},
			startTime: "09:00",
			endTime:   "10:00",
			setup: func(roomRepo *mocks.MockRoomRepository, scheduleRepo *mocks.MockScheduleRepository, slotRepo *mocks.MockSlotRepository) {
				roomRepo.EXPECT().Exists(mock.Anything, "room-1").Return(true, nil)
				scheduleRepo.EXPECT().Create(mock.Anything, "room-1", []int{1}, "09:00", "10:00").
					Return(nil, storage.ErrInvalidTimeRange)
			},
			wantErr: domain.ErrInvalidTimeRange,
		},
		{
			name:      "storage error",
			roomID:    "room-1",
			weekdays:  []int{1},
			startTime: "09:00",
			endTime:   "10:00",
			setup: func(roomRepo *mocks.MockRoomRepository, scheduleRepo *mocks.MockScheduleRepository, slotRepo *mocks.MockSlotRepository) {
				roomRepo.EXPECT().Exists(mock.Anything, "room-1").Return(true, nil)
				scheduleRepo.EXPECT().Create(mock.Anything, "room-1", []int{1}, "09:00", "10:00").
					Return(nil, errDummy)
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

			uc := New(scheduleRepo, roomRepo, slotRepo)

			schedule, err := uc.CreateSchedule(t.Context(), tt.roomID, tt.weekdays, tt.startTime, tt.endTime)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, schedule)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, schedule)
		})
	}
}

func testUUID() uuid.UUID {
	return uuid.MustParse("00000000-0000-0000-0000-000000000011")
}
