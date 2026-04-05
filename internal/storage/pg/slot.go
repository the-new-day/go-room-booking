package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain/entity"
	"github.com/internships-backend/test-backend-the-new-day/internal/storage"
	"github.com/internships-backend/test-backend-the-new-day/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

type SlotRepository struct {
	db *postgres.Postgres
}

func NewSlotRepository(db *postgres.Postgres) *SlotRepository {
	return &SlotRepository{db: db}
}

func (r *SlotRepository) Create(ctx context.Context, roomID string, startAt, endAt time.Time) error {
	const op = "storage.pg.SlotRepository.Create"

	query, args, err := r.db.Builder.
		Insert("slots").
		Columns("room_id", "start_time", "end_time").
		Values(roomID, startAt, endAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = r.db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, handlePostgresError(err))
	}

	return nil
}

func (r *SlotRepository) ListAvailableByRoomAndDate(ctx context.Context, roomID string, date time.Time) ([]*entity.Slot, error) {
	const op = "storage.pg.SlotRepository.ListAvailableByRoomAndDate"

	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	query, args, err := r.db.Builder.
		Select("id", "room_id", "start_time", "end_time").
		From("slots").
		Where(squirrel.Eq{"room_id": roomID, "is_available": true}).
		Where(squirrel.GtOrEq{"start_time": startOfDay}).
		Where(squirrel.Lt{"start_time": endOfDay}).
		OrderBy("start_time ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var slots []*entity.Slot
	for rows.Next() {
		var slot entity.Slot
		err := rows.Scan(&slot.SlotID, &slot.RoomID, &slot.StartAt, &slot.EndAt)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		slots = append(slots, &slot)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return slots, nil
}

func (r *SlotRepository) GetByID(ctx context.Context, slotID string) (*entity.Slot, error) {
	const op = "storage.pg.SlotRepository.GetByID"

	query, args, err := r.db.Builder.
		Select("id", "room_id", "start_time", "end_time").
		From("slots").
		Where(squirrel.Eq{"id": slotID}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var slot entity.Slot
	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(&slot.SlotID, &slot.RoomID, &slot.StartAt, &slot.EndAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, handlePostgresError(err))
	}

	return &slot, nil
}
