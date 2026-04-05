package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain/entity"
	"github.com/internships-backend/test-backend-the-new-day/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

type ScheduleRepository struct {
	db *postgres.Postgres
}

func NewScheduleRepository(db *postgres.Postgres) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

func (r *ScheduleRepository) ExistsByRoomID(ctx context.Context, roomID string) (bool, error) {
	const op = "storage.pg.ScheduleRepository.ExistsByRoomID"

	query, args, err := r.db.Builder.
		Select("1").
		From("schedules").
		Where(squirrel.Eq{"room_id": roomID}).
		Limit(1).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	var exists int
	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return true, nil
}

func (r *ScheduleRepository) Create(ctx context.Context, roomID string, weekdays []int, startTime, endTime string) (*entity.Schedule, error) {
	const op = "storage.pg.ScheduleRepository.Create"

	query, args, err := r.db.Builder.
		Insert("schedules").
		Columns("room_id", "days_of_week", "start_time", "end_time").
		Values(roomID, weekdays, startTime, endTime).
		Suffix("RETURNING id, room_id, days_of_week, start_time, end_time, created_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var schedule entity.Schedule
	var days []int

	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(
		&schedule.ScheduleID,
		&schedule.RoomID,
		&days,
		&schedule.StartAt,
		&schedule.EndAt,
		&schedule.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, handlePostgresError(err))
	}

	schedule.Weekdays = make([]entity.Weekday, len(days))
	for i, day := range days {
		schedule.Weekdays[i] = entity.Weekday(day)
	}

	return &schedule, nil
}
