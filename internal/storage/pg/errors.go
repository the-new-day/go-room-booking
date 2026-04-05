package pg

import (
	"errors"
	"strings"

	"github.com/internships-backend/test-backend-the-new-day/internal/storage"
	"github.com/jackc/pgx/v5/pgconn"
)

func handlePostgresError(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return handleUniqueViolation(pgErr)
		case "23514": // check_violation
			return handleCheckViolation(pgErr)
		case "23P01": // exclusion_violation
			if strings.Contains(pgErr.ConstraintName, "unique_room_time_range") {
				return storage.ErrSlotOverlap
			}
			return storage.ErrAlreadyExists
		}
	}

	return err
}

func handleUniqueViolation(pgErr *pgconn.PgError) error {
	switch pgErr.ConstraintName {
	case "users_email_key":
		return storage.ErrDuplicateEmail
	case "schedules_room_id_key":
		return storage.ErrScheduleAlreadyExists
	case "unique_active_booking_per_slot":
		return storage.ErrSlotAlreadyBooked
	default:
		return storage.ErrAlreadyExists
	}
}

func handleCheckViolation(pgErr *pgconn.PgError) error {
	switch pgErr.ConstraintName {
	case "users_role_check":
		return storage.ErrInvalidRole
	case "rooms_capacity_check":
		return storage.ErrInvalidCapacity
	case "schedules_valid_days_of_week_check":
		return storage.ErrInvalidDaysOfWeek
	case "schedules_valid_time_range_check", "slots_valid_time_range_check":
		return storage.ErrInvalidTimeRange
	}
	return pgErr
}
