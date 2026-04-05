package schedule

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/internships-backend/test-backend-the-new-day/internal/domain"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain/entity"
	"github.com/internships-backend/test-backend-the-new-day/internal/storage"
)

type ScheduleRepository interface {
	ExistsByRoomID(ctx context.Context, roomID string) (bool, error)
	Create(ctx context.Context, roomID string, weekdays []int, startTime, endTime string) (*entity.Schedule, error)
}

type RoomRepository interface {
	Exists(ctx context.Context, roomID string) (bool, error)
}

type SlotRepository interface {
	Create(ctx context.Context, roomID string, startAt, endAt time.Time) error
}

type UseCase struct {
	scheduleRepo ScheduleRepository
	roopRepo     RoomRepository
	slotRepo     SlotRepository
}

func New(scheduleRepo ScheduleRepository, roomRepo RoomRepository, slotRepo SlotRepository) *UseCase {
	return &UseCase{
		scheduleRepo: scheduleRepo,
		roopRepo:     roomRepo,
		slotRepo:     slotRepo,
	}
}

func (u *UseCase) CreateSchedule(ctx context.Context, roomID string, weekdays []int, startTime, endTime string) (*entity.Schedule, error) {
	const op = "usecase.schedule.CreateSchedule"

	exists, err := u.roopRepo.Exists(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if !exists {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrRoomNotFound)
	}

	if len(weekdays) == 0 {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrInvalidDaysOfWeek)
	}

	for _, day := range weekdays {
		if day < 1 || day > 7 {
			return nil, fmt.Errorf("%s: %w", op, domain.ErrInvalidDaysOfWeek)
		}
	}

	startParsed, err := time.Parse("15:04", startTime)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrInvalidTimeRange)
	}

	endParsed, err := time.Parse("15:04", endTime)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrInvalidTimeRange)
	}

	if !startParsed.Before(endParsed) {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrInvalidTimeRange)
	}

	schedule, err := u.scheduleRepo.Create(ctx, roomID, weekdays, startTime, endTime)
	if err != nil {
		if errors.Is(err, storage.ErrScheduleAlreadyExists) {
			return nil, fmt.Errorf("%s: %w", op, domain.ErrScheduleAlreadyExists)
		}
		if errors.Is(err, storage.ErrInvalidDaysOfWeek) {
			return nil, fmt.Errorf("%s: %w", op, domain.ErrInvalidDaysOfWeek)
		}
		if errors.Is(err, storage.ErrInvalidTimeRange) {
			return nil, fmt.Errorf("%s: %w", op, domain.ErrInvalidTimeRange)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = u.generateSlots(ctx, roomID, weekdays, startParsed, endParsed)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return schedule, nil
}

func (u *UseCase) generateSlots(ctx context.Context, roomID string, weekdays []int, startTime, endTime time.Time) error {
	startDate := time.Now().UTC()
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)

	for i := 0; i <= 7; i++ {
		day := startDate.AddDate(0, 0, i)
		if !containsWeekday(weekdays, day.Weekday()) {
			continue
		}

		slotStart := time.Date(day.Year(), day.Month(), day.Day(), startTime.Hour(), startTime.Minute(), 0, 0, time.UTC)
		slotEnd := time.Date(day.Year(), day.Month(), day.Day(), endTime.Hour(), endTime.Minute(), 0, 0, time.UTC)

		for t := slotStart; t.Add(entity.SlotDuration).Equal(slotEnd) || t.Add(entity.SlotDuration).Before(slotEnd); t = t.Add(entity.SlotDuration) {
			if err := u.slotRepo.Create(ctx, roomID, t, t.Add(entity.SlotDuration)); err != nil {
				return err
			}
		}
	}

	return nil
}

func containsWeekday(weekdays []int, day time.Weekday) bool {
	target := int(day)
	if target == 0 {
		target = 7
	}
	for _, w := range weekdays {
		if w == target {
			return true
		}
	}
	return false
}
