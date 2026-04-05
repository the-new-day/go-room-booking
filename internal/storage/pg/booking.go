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

type BookingRepository struct {
	db *postgres.Postgres
}

func NewBookingRepository(db *postgres.Postgres) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(ctx context.Context, slotID, userID string, conferenceLink *string) (*entity.Booking, error) {
	const op = "storage.pg.BookingRepository.Create"

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback(ctx)

	updateQuery, updateArgs, err := r.db.Builder.
		Update("slots").
		Set("is_available", false).
		Where(squirrel.Eq{"id": slotID, "is_available": true}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ct, err := tx.Exec(ctx, updateQuery, updateArgs...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, handlePostgresError(err))
	}
	if ct.RowsAffected() == 0 {
		return nil, fmt.Errorf("%s: %w", op, storage.ErrSlotAlreadyBooked)
	}

	insertQuery, insertArgs, err := r.db.Builder.
		Insert("bookings").
		Columns("slot_id", "user_id", "status", "conference_link").
		Values(slotID, userID, entity.BookingActive, conferenceLink).
		Suffix("RETURNING id, slot_id, user_id, status, conference_link, created_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var booking entity.Booking
	err = tx.QueryRow(ctx, insertQuery, insertArgs...).Scan(
		&booking.BookingID,
		&booking.SlotID,
		&booking.UserID,
		&booking.Status,
		&booking.ConferenceLink,
		&booking.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, handlePostgresError(err))
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &booking, nil
}

func (r *BookingRepository) FindByID(ctx context.Context, bookingID string) (*entity.Booking, error) {
	const op = "storage.pg.BookingRepository.FindByID"

	query, args, err := r.db.Builder.
		Select("id", "slot_id", "user_id", "status", "conference_link", "created_at").
		From("bookings").
		Where(squirrel.Eq{"id": bookingID}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var booking entity.Booking
	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(
		&booking.BookingID,
		&booking.SlotID,
		&booking.UserID,
		&booking.Status,
		&booking.ConferenceLink,
		&booking.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, handlePostgresError(err))
	}

	return &booking, nil
}

func (r *BookingRepository) Cancel(ctx context.Context, bookingID string) (*entity.Booking, error) {
	const op = "storage.pg.BookingRepository.Cancel"

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback(ctx)

	query, args, err := r.db.Builder.
		Update("bookings").
		Set("status", entity.BookingCancelled).
		Where(squirrel.Eq{"id": bookingID, "status": entity.BookingActive}).
		Suffix("RETURNING id, slot_id, user_id, status, conference_link, created_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var booking entity.Booking
	err = tx.QueryRow(ctx, query, args...).Scan(
		&booking.BookingID,
		&booking.SlotID,
		&booking.UserID,
		&booking.Status,
		&booking.ConferenceLink,
		&booking.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, handlePostgresError(err))
	}

	updateSlotQuery, updateSlotArgs, err := r.db.Builder.
		Update("slots").
		Set("is_available", true).
		Where(squirrel.Eq{"id": booking.SlotID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.Exec(ctx, updateSlotQuery, updateSlotArgs...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, handlePostgresError(err))
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &booking, nil
}

func (r *BookingRepository) List(ctx context.Context, offset, limit int) ([]*entity.Booking, int, error) {
	const op = "storage.pg.BookingRepository.List"

	var total int
	err := r.db.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM bookings").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	query, args, err := r.db.Builder.
		Select("id", "slot_id", "user_id", "status", "conference_link", "created_at").
		From("bookings").
		OrderBy("created_at DESC").
		Offset(uint64(offset)).
		Limit(uint64(limit)).
		ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var bookings []*entity.Booking
	for rows.Next() {
		var booking entity.Booking
		err := rows.Scan(
			&booking.BookingID,
			&booking.SlotID,
			&booking.UserID,
			&booking.Status,
			&booking.ConferenceLink,
			&booking.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("%s: %w", op, err)
		}
		bookings = append(bookings, &booking)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	return bookings, total, nil
}

func (r *BookingRepository) ListByUserFuture(ctx context.Context, userID string, now time.Time) ([]*entity.Booking, error) {
	const op = "storage.pg.BookingRepository.ListByUserFuture"

	query, args, err := r.db.Builder.
		Select("b.id", "b.slot_id", "b.user_id", "b.status", "b.conference_link", "b.created_at").
		From("bookings b").
		Join("slots s ON s.id = b.slot_id").
		Where(squirrel.Eq{"b.user_id": userID}).
		Where(squirrel.GtOrEq{"s.start_time": now}).
		OrderBy("s.start_time ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var bookings []*entity.Booking
	for rows.Next() {
		var booking entity.Booking
		err := rows.Scan(
			&booking.BookingID,
			&booking.SlotID,
			&booking.UserID,
			&booking.Status,
			&booking.ConferenceLink,
			&booking.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		bookings = append(bookings, &booking)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return bookings, nil
}
