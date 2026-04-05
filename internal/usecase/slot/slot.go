package slot

import (
	"context"
	"fmt"
	"time"

	"github.com/internships-backend/test-backend-the-new-day/internal/domain"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain/entity"
)

type RoomRepository interface {
	Exists(ctx context.Context, roomID string) (bool, error)
}

type ScheduleRepository interface {
	ExistsByRoomID(ctx context.Context, roomID string) (bool, error)
}

type SlotRepository interface {
	ListAvailableByRoomAndDate(ctx context.Context, roomID string, date time.Time) ([]*entity.Slot, error)
}

type UseCase struct {
	roomRepo     RoomRepository
	scheduleRepo ScheduleRepository
	slotRepo     SlotRepository
}

func New(roomRepo RoomRepository, scheduleRepo ScheduleRepository, slotRepo SlotRepository) *UseCase {
	return &UseCase{
		roomRepo:     roomRepo,
		scheduleRepo: scheduleRepo,
		slotRepo:     slotRepo,
	}
}

func (u *UseCase) ListAvailableSlots(ctx context.Context, roomID string, date time.Time) ([]*entity.Slot, error) {
	const op = "usecase.slot.ListAvailableSlots"

	exists, err := u.roomRepo.Exists(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if !exists {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrRoomNotFound)
	}

	hasSchedule, err := u.scheduleRepo.ExistsByRoomID(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if !hasSchedule {
		return []*entity.Slot{}, nil
	}

	return u.slotRepo.ListAvailableByRoomAndDate(ctx, roomID, date)
}
