package booking

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain/entity"
	"github.com/internships-backend/test-backend-the-new-day/internal/storage"
	"github.com/internships-backend/test-backend-the-new-day/internal/usecase/booking/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBookingUseCase_CreateBooking(t *testing.T) {
	slotID := uuid.New().String()
	userID := uuid.New().String()
	now := time.Date(2026, 4, 5, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		setup   func(bookingRepo *mocks.MockBookingRepository, slotRepo *mocks.MockSlotRepository)
		wantErr error
	}{
		{
			name: "slot not found",
			setup: func(bookingRepo *mocks.MockBookingRepository, slotRepo *mocks.MockSlotRepository) {
				slotRepo.EXPECT().GetByID(mock.Anything, slotID).Return(nil, storage.ErrNotFound)
			},
			wantErr: domain.ErrSlotNotFound,
		},
		{
			name: "slot in past",
			setup: func(bookingRepo *mocks.MockBookingRepository, slotRepo *mocks.MockSlotRepository) {
				slotRepo.EXPECT().GetByID(mock.Anything, slotID).
					Return(&entity.Slot{StartAt: now.Add(-time.Minute)}, nil)
			},
			wantErr: domain.ErrSlotInPast,
		},
		{
			name: "slot already booked",
			setup: func(bookingRepo *mocks.MockBookingRepository, slotRepo *mocks.MockSlotRepository) {
				slotRepo.EXPECT().GetByID(mock.Anything, slotID).
					Return(&entity.Slot{StartAt: now.Add(time.Minute)}, nil)
				bookingRepo.EXPECT().Create(mock.Anything, slotID, userID, (*string)(nil)).
					Return(nil, storage.ErrSlotAlreadyBooked)
			},
			wantErr: domain.ErrSlotAlreadyBooked,
		},
		{
			name: "success",
			setup: func(bookingRepo *mocks.MockBookingRepository, slotRepo *mocks.MockSlotRepository) {
				slotRepo.EXPECT().GetByID(mock.Anything, slotID).
					Return(&entity.Slot{StartAt: now.Add(time.Minute)}, nil)
				bookingRepo.EXPECT().Create(mock.Anything, slotID, userID, (*string)(nil)).
					Return(&entity.Booking{BookingID: uuid.New()}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			bookingRepo := mocks.NewMockBookingRepository(t)
			slotRepo := mocks.NewMockSlotRepository(t)

			if tt.setup != nil {
				tt.setup(bookingRepo, slotRepo)
			}

			uc := New(bookingRepo, slotRepo)
			uc.now = func() time.Time { return now }

			booking, err := uc.CreateBooking(t.Context(), slotID, userID, false)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, booking)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, booking)
		})
	}
}

func TestBookingUseCase_CancelBooking(t *testing.T) {
	bookingID := uuid.New().String()
	userID := uuid.New().String()
	otherUserID := uuid.New().String()

	tests := []struct {
		name    string
		setup   func(bookingRepo *mocks.MockBookingRepository)
		wantErr error
	}{
		{
			name: "booking not found",
			setup: func(bookingRepo *mocks.MockBookingRepository) {
				bookingRepo.EXPECT().FindByID(mock.Anything, bookingID).Return(nil, storage.ErrNotFound)
			},
			wantErr: domain.ErrBookingNotFound,
		},
		{
			name: "forbidden",
			setup: func(bookingRepo *mocks.MockBookingRepository) {
				bookingRepo.EXPECT().FindByID(mock.Anything, bookingID).
					Return(&entity.Booking{UserID: uuid.MustParse(otherUserID), Status: entity.BookingActive}, nil)
			},
			wantErr: domain.ErrForbidden,
		},
		{
			name: "already cancelled",
			setup: func(bookingRepo *mocks.MockBookingRepository) {
				bookingRepo.EXPECT().FindByID(mock.Anything, bookingID).
					Return(&entity.Booking{UserID: uuid.MustParse(userID), Status: entity.BookingCancelled}, nil)
			},
		},
		{
			name: "success",
			setup: func(bookingRepo *mocks.MockBookingRepository) {
				bookingRepo.EXPECT().FindByID(mock.Anything, bookingID).
					Return(&entity.Booking{UserID: uuid.MustParse(userID), Status: entity.BookingActive}, nil)
				bookingRepo.EXPECT().Cancel(mock.Anything, bookingID).
					Return(&entity.Booking{UserID: uuid.MustParse(userID), Status: entity.BookingCancelled}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			bookingRepo := mocks.NewMockBookingRepository(t)
			slotRepo := mocks.NewMockSlotRepository(t)

			if tt.setup != nil {
				tt.setup(bookingRepo)
			}

			uc := New(bookingRepo, slotRepo)

			booking, err := uc.CancelBooking(t.Context(), bookingID, userID)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, booking)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, booking)
		})
	}
}

func TestBookingUseCase_ListBookings(t *testing.T) {
	bookingRepo := mocks.NewMockBookingRepository(t)
	slotRepo := mocks.NewMockSlotRepository(t)

	bookingRepo.EXPECT().List(mock.Anything, 20, 10).Return([]*entity.Booking{{}}, 1, nil)

	uc := New(bookingRepo, slotRepo)

	bookings, total, err := uc.ListBookings(t.Context(), 3, 10)
	require.NoError(t, err)
	assert.Len(t, bookings, 1)
	assert.Equal(t, 1, total)
}

func TestBookingUseCase_ListMyBookings(t *testing.T) {
	bookingRepo := mocks.NewMockBookingRepository(t)
	slotRepo := mocks.NewMockSlotRepository(t)

	bookingRepo.EXPECT().ListByUserFuture(mock.Anything, "user-1", mock.Anything).Return([]*entity.Booking{{}}, nil)

	uc := New(bookingRepo, slotRepo)
	uc.now = func() time.Time { return time.Date(2026, 4, 5, 10, 0, 0, 0, time.UTC) }

	bookings, err := uc.ListMyBookings(t.Context(), "user-1")
	require.NoError(t, err)
	assert.Len(t, bookings, 1)
}
