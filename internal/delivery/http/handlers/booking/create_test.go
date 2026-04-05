package booking

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/handlers/booking/mocks"
	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/middleware"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain/entity"
	"github.com/internships-backend/test-backend-the-new-day/pkg/logger/sl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBookingCreateHandler(t *testing.T) {
	userID := uuid.New().String()
	slotID := uuid.New().String()

	tests := []struct {
		name       string
		body       any
		withUser   bool
		setup      func(creator *mocks.MockBookingCreator)
		wantStatus int
	}{
		{
			name:       "invalid body",
			body:       map[string]any{},
			withUser:   true,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid slot id",
			body:       map[string]any{"slotId": "invalid"},
			withUser:   true,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "no user",
			body:       map[string]any{"slotId": slotID},
			withUser:   false,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:     "slot not found",
			body:     map[string]any{"slotId": slotID},
			withUser: true,
			setup: func(creator *mocks.MockBookingCreator) {
				creator.EXPECT().CreateBooking(mock.Anything, slotID, userID, false).
					Return(nil, domain.ErrSlotNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:     "slot already booked",
			body:     map[string]any{"slotId": slotID},
			withUser: true,
			setup: func(creator *mocks.MockBookingCreator) {
				creator.EXPECT().CreateBooking(mock.Anything, slotID, userID, false).
					Return(nil, domain.ErrSlotAlreadyBooked)
			},
			wantStatus: http.StatusConflict,
		},
		{
			name:     "slot in past",
			body:     map[string]any{"slotId": slotID},
			withUser: true,
			setup: func(creator *mocks.MockBookingCreator) {
				creator.EXPECT().CreateBooking(mock.Anything, slotID, userID, false).
					Return(nil, domain.ErrSlotInPast)
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:     "success",
			body:     map[string]any{"slotId": slotID},
			withUser: true,
			setup: func(creator *mocks.MockBookingCreator) {
				creator.EXPECT().CreateBooking(mock.Anything, slotID, userID, false).
					Return(&entity.Booking{BookingID: uuid.New()}, nil)
			},
			wantStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			creator := mocks.NewMockBookingCreator(t)
			if tt.setup != nil {
				tt.setup(creator)
			}

			handler := NewCreateHandler(sl.NewDiscardLogger(), creator)

			body := bytes.NewBuffer(nil)
			if tt.body != nil {
				require.NoError(t, json.NewEncoder(body).Encode(tt.body))
			}

			req := httptest.NewRequest(http.MethodPost, "/bookings/create", body)
			req.Header.Set("Content-Type", "application/json")

			if tt.withUser {
				req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, userID))
			}

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
