package room

import (
	"context"
	"strings"

	"github.com/internships-backend/test-backend-the-new-day/internal/domain"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain/entity"
)

type RoomRepository interface {
	Create(ctx context.Context, name string, description *string, capacity *int) (*entity.Room, error)
	List(ctx context.Context) ([]*entity.Room, error)
}

type UseCase struct {
	repo RoomRepository
}

func New(repo RoomRepository) *UseCase {
	return &UseCase{
		repo: repo,
	}
}

func (u *UseCase) CreateRoom(ctx context.Context, name string, description *string, capacity *int) (*entity.Room, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, domain.ErrEmptyRoomName
	}

	if capacity != nil && *capacity <= 0 {
		return nil, domain.ErrNonPositiveRoomCapacity
	}

	return u.repo.Create(ctx, name, description, capacity)
}

func (u *UseCase) ListRooms(ctx context.Context) ([]*entity.Room, error) {
	return u.repo.List(ctx)
}
