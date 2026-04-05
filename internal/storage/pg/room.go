package pg

import (
	"context"
	"fmt"

	"github.com/internships-backend/test-backend-the-new-day/internal/domain/entity"
	"github.com/internships-backend/test-backend-the-new-day/pkg/postgres"
)

type RoomRepository struct {
	db *postgres.Postgres
}

func NewRoomRepository(db *postgres.Postgres) *RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) Create(ctx context.Context, name string, description *string, capacity *int) (*entity.Room, error) {
	const op = "storage.pg.RoomRepository.Create"

	query, args, err := r.db.Builder.
		Insert("rooms").
		Columns("name", "description", "capacity").
		Values(name, description, capacity).
		Suffix("RETURNING id, name, description, capacity, created_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var room entity.Room

	row := r.db.Pool.QueryRow(ctx, query, args...)
	err = row.Scan(&room.RoomID, &room.Name, &room.Description, &room.Capacity, &room.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, handlePostgresError(err))
	}

	return &room, nil
}

func (r *RoomRepository) List(ctx context.Context) ([]*entity.Room, error) {
	const op = "storage.pg.RoomRepository.List"

	query, args, err := r.db.Builder.
		Select("id", "name", "description", "capacity", "created_at").
		From("rooms").
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var rooms []*entity.Room

	for rows.Next() {
		var room entity.Room
		err := rows.Scan(
			&room.RoomID,
			&room.Name,
			&room.Description,
			&room.Capacity,
			&room.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		rooms = append(rooms, &room)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return rooms, nil
}
