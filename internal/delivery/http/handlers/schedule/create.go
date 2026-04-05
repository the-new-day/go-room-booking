package schedule

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/api"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain/entity"
	"github.com/internships-backend/test-backend-the-new-day/pkg/logger/sl"
)

type CreateRequest struct {
	RoomID     string `json:"roomId" validate:"required"`
	DaysOfWeek []int  `json:"daysOfWeek" validate:"required,min=1,dive,min=1,max=7"`
	StartTime  string `json:"startTime" validate:"required"`
	EndTime    string `json:"endTime" validate:"required"`
}

type ScheduleResponse struct {
	ScheduleID string `json:"id"`
	RoomID     string `json:"roomId"`
	DaysOfWeek []int  `json:"daysOfWeek"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
}

type CreateResponse struct {
	Schedule ScheduleResponse `json:"schedule"`
}

type ScheduleCreator interface {
	CreateSchedule(ctx context.Context, roomID string, weekdays []int, startTime, endTime string) (*entity.Schedule, error)
}

func NewCreateHandler(logger *slog.Logger, scheduleCreator ScheduleCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomID := chi.URLParam(r, "roomId")
		if _, err := uuid.Parse(roomID); err != nil {
			api.SendBadRequest(w, r, "invalid room id")
			return
		}

		req, ok := api.DecodeRequest[CreateRequest](logger, w, r)
		if !ok {
			return
		}

		if _, err := uuid.Parse(req.RoomID); err != nil {
			api.SendBadRequest(w, r, "invalid room id")
			return
		}

		if req.RoomID != roomID {
			api.SendBadRequest(w, r, "room id mismatch")
			return
		}

		schedule, err := scheduleCreator.CreateSchedule(r.Context(), roomID, req.DaysOfWeek, req.StartTime, req.EndTime)
		if errors.Is(err, domain.ErrRoomNotFound) {
			api.SendRoomNotFound(w, r, "room not found")
			return
		} else if errors.Is(err, domain.ErrScheduleAlreadyExists) {
			api.SendScheduleExists(w, r, "schedule for this room already exists and cannot be changed")
			return
		} else if errors.Is(err, domain.ErrInvalidDaysOfWeek) || errors.Is(err, domain.ErrInvalidTimeRange) {
			api.SendBadRequest(w, r, "invalid schedule")
			return
		} else if err != nil {
			logger.Error("failed to create schedule", sl.Err(err))
			api.SendInternalError(w, r, "failed to create schedule")
			return
		}

		api.SendCreated(w, r, CreateResponse{
			Schedule: ScheduleResponse{
				ScheduleID: schedule.ScheduleID.String(),
				RoomID:     schedule.RoomID.String(),
				DaysOfWeek: mapWeekdays(schedule.Weekdays),
				StartTime:  schedule.StartAt.Format("15:04"),
				EndTime:    schedule.EndAt.Format("15:04"),
			},
		})
	}
}

func mapWeekdays(days []entity.Weekday) []int {
	res := make([]int, len(days))
	for i, day := range days {
		res[i] = int(day)
	}
	return res
}
