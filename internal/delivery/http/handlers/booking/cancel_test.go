package booking

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/handlers/booking/mocks"
	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/middleware"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain/entity"
	"github.com/internships-backend/test-backend-the-new-day/pkg/logger/sl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBookingCancelHandler(t *testing.T) {
	bookingID := uuid.New().String()

	tests := []struct {
		name       string
		bookingID  string
		withUser   bool
		setup      func(canceler *mocks.MockBookingCanceler)
		wantStatus int
	}{
		{
			name:       "invalid booking id",
			bookingID:  "invalid",
			withUser:   true,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "no user",
			bookingID:  bookingID,
			withUser:   false,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:      "booking not found",
			bookingID: bookingID,
			withUser:  true,
			setup: func(canceler *mocks.MockBookingCanceler) {
				canceler.EXPECT().CancelBooking(mock.Anything, bookingID, "user-1").
					Return(nil, domain.ErrBookingNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:      "forbidden",
			bookingID: bookingID,
			withUser:  true,
			setup: func(canceler *mocks.MockBookingCanceler) {
				canceler.EXPECT().CancelBooking(mock.Anything, bookingID, "user-1").
					Return(nil, domain.ErrForbidden)
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name:      "success",
			bookingID: bookingID,
			withUser:  true,
			setup: func(canceler *mocks.MockBookingCanceler) {
				canceler.EXPECT().CancelBooking(mock.Anything, bookingID, "user-1").
					Return(&entity.Booking{BookingID: uuid.MustParse(bookingID)}, nil)
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			canceler := mocks.NewMockBookingCanceler(t)
			if tt.setup != nil {
				tt.setup(canceler)
			}

			handler := NewCancelHandler(sl.NewDiscardLogger(), canceler)

			req := httptest.NewRequest(http.MethodPost, "/bookings/"+tt.bookingID+"/cancel", nil)
			if tt.withUser {
				req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, "user-1"))
			}
			req = withURLParam(req, "bookingId", tt.bookingID)

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func withURLParam(r *http.Request, key, value string) *http.Request {
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add(key, value)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, routeCtx))
}
