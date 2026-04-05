package schedule

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/handlers/schedule/mocks"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain/entity"
	"github.com/internships-backend/test-backend-the-new-day/pkg/logger/sl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestScheduleCreateHandler(t *testing.T) {
	roomID := uuid.New().String()

	tests := []struct {
		name       string
		roomID     string
		body       any
		setup      func(creator *mocks.MockScheduleCreator)
		wantStatus int
	}{
		{
			name:       "invalid room id path",
			roomID:     "invalid",
			body:       map[string]any{"roomId": roomID, "daysOfWeek": []int{1}, "startTime": "09:00", "endTime": "10:00"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid body",
			roomID:     roomID,
			body:       map[string]any{"roomId": roomID},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "room id mismatch",
			roomID:     roomID,
			body:       map[string]any{"roomId": uuid.New().String(), "daysOfWeek": []int{1}, "startTime": "09:00", "endTime": "10:00"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:   "room not found",
			roomID: roomID,
			body:   map[string]any{"roomId": roomID, "daysOfWeek": []int{1}, "startTime": "09:00", "endTime": "10:00"},
			setup: func(creator *mocks.MockScheduleCreator) {
				creator.EXPECT().CreateSchedule(mock.Anything, roomID, []int{1}, "09:00", "10:00").
					Return(nil, domain.ErrRoomNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:   "schedule exists",
			roomID: roomID,
			body:   map[string]any{"roomId": roomID, "daysOfWeek": []int{1}, "startTime": "09:00", "endTime": "10:00"},
			setup: func(creator *mocks.MockScheduleCreator) {
				creator.EXPECT().CreateSchedule(mock.Anything, roomID, []int{1}, "09:00", "10:00").
					Return(nil, domain.ErrScheduleAlreadyExists)
			},
			wantStatus: http.StatusConflict,
		},
		{
			name:   "success",
			roomID: roomID,
			body:   map[string]any{"roomId": roomID, "daysOfWeek": []int{1}, "startTime": "09:00", "endTime": "10:00"},
			setup: func(creator *mocks.MockScheduleCreator) {
				creator.EXPECT().CreateSchedule(mock.Anything, roomID, []int{1}, "09:00", "10:00").
					Return(&entity.Schedule{
						ScheduleID: uuid.MustParse("00000000-0000-0000-0000-000000000021"),
						RoomID:     uuid.MustParse(roomID),
						Weekdays:   []entity.Weekday{entity.Monday},
						StartAt:    time.Date(2026, 4, 5, 9, 0, 0, 0, time.UTC),
						EndAt:      time.Date(2026, 4, 5, 10, 0, 0, 0, time.UTC),
					}, nil)
			},
			wantStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			creator := mocks.NewMockScheduleCreator(t)
			if tt.setup != nil {
				tt.setup(creator)
			}

			handler := NewCreateHandler(sl.NewDiscardLogger(), creator)

			body := bytes.NewBuffer(nil)
			if tt.body != nil {
				require.NoError(t, json.NewEncoder(body).Encode(tt.body))
			}

			req := httptest.NewRequest(http.MethodPost, "/rooms/"+tt.roomID+"/schedule/create", body)
			req.Header.Set("Content-Type", "application/json")
			req = withURLParam(req, "roomId", tt.roomID)

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
