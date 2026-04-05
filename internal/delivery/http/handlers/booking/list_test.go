package booking

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/handlers/booking/mocks"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain/entity"
	"github.com/internships-backend/test-backend-the-new-day/pkg/logger/sl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBookingListHandler(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		setup      func(lister *mocks.MockBookingLister)
		wantStatus int
	}{
		{
			name:       "invalid page",
			query:      "?page=0",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid page size",
			query:      "?pageSize=101",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:  "success",
			query: "?page=2&pageSize=10",
			setup: func(lister *mocks.MockBookingLister) {
				lister.EXPECT().ListBookings(mock.Anything, 2, 10).Return([]*entity.Booking{{}}, 1, nil)
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			lister := mocks.NewMockBookingLister(t)
			if tt.setup != nil {
				tt.setup(lister)
			}

			handler := NewListHandler(sl.NewDiscardLogger(), lister)

			req := httptest.NewRequest(http.MethodGet, "/bookings/list"+tt.query, nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
