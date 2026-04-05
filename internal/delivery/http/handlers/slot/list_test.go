package slot

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/handlers/slot/mocks"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain/entity"
	"github.com/internships-backend/test-backend-the-new-day/pkg/logger/sl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSlotListHandler(t *testing.T) {
	roomID := uuid.New().String()

	tests := []struct {
		name       string
		roomID     string
		date       string
		setup      func(lister *mocks.MockSlotLister)
		wantStatus int
	}{
		{
			name:       "invalid room id",
			roomID:     "invalid",
			date:       "2026-04-05",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing date",
			roomID:     roomID,
			date:       "",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid date",
			roomID:     roomID,
			date:       "2026-13-01",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:   "room not found",
			roomID: roomID,
			date:   "2026-04-05",
			setup: func(lister *mocks.MockSlotLister) {
				lister.EXPECT().ListAvailableSlots(mock.Anything, roomID, time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC)).
					Return(nil, domain.ErrRoomNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:   "success",
			roomID: roomID,
			date:   "2026-04-05",
			setup: func(lister *mocks.MockSlotLister) {
				lister.EXPECT().ListAvailableSlots(mock.Anything, roomID, time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC)).
					Return([]*entity.Slot{{SlotID: uuid.MustParse("00000000-0000-0000-0000-000000000031")}}, nil)
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			lister := mocks.NewMockSlotLister(t)
			if tt.setup != nil {
				tt.setup(lister)
			}

			handler := NewListHandler(sl.NewDiscardLogger(), lister)

			req := httptest.NewRequest(http.MethodGet, "/rooms/"+tt.roomID+"/slots/list", nil)
			q := req.URL.Query()
			if tt.date != "" {
				q.Set("date", tt.date)
			}
			req.URL.RawQuery = q.Encode()
			req = withURLParam(req, "roomId", tt.roomID)

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			if rec.Code == http.StatusOK {
				var resp ListResponse
				_ = json.NewDecoder(rec.Body).Decode(&resp)
				assert.NotNil(t, resp.Slots)
			}
		})
	}
}

func withURLParam(r *http.Request, key, value string) *http.Request {
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add(key, value)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, routeCtx))
}
