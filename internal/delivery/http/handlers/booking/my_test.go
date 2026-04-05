package booking

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/handlers/booking/mocks"
	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/middleware"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain/entity"
	"github.com/internships-backend/test-backend-the-new-day/pkg/logger/sl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBookingMyHandler(t *testing.T) {
	tests := []struct {
		name       string
		withUser   bool
		setup      func(lister *mocks.MockMyBookingLister)
		wantStatus int
	}{
		{
			name:       "no user",
			withUser:   false,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:     "success",
			withUser: true,
			setup: func(lister *mocks.MockMyBookingLister) {
				lister.EXPECT().ListMyBookings(mock.Anything, "user-1").Return([]*entity.Booking{{}}, nil)
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			lister := mocks.NewMockMyBookingLister(t)
			if tt.setup != nil {
				tt.setup(lister)
			}

			handler := NewMyHandler(sl.NewDiscardLogger(), lister)

			req := httptest.NewRequest(http.MethodGet, "/bookings/my", nil)
			if tt.withUser {
				req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, "user-1"))
			}
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
