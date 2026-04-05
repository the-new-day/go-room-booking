package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain/entity"
	"github.com/internships-backend/test-backend-the-new-day/internal/storage"
	"github.com/internships-backend/test-backend-the-new-day/pkg/postgres"
)

type UserRepository struct {
	db *postgres.Postgres
}

func NewUserRepository(db *postgres.Postgres) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, email, passwordHash, role string) (*entity.User, error) {
	const op = "storage.pg.UserRepository.Create"

	query, args, err := r.db.Builder.
		Insert("users").
		Columns("email", "password_hash", "role").
		Values(email, passwordHash, role).
		Suffix("RETURNING id, email, role, created_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var user entity.User

	row := r.db.Pool.QueryRow(ctx, query, args...)
	err = row.Scan(&user.UserID, &user.Email, &user.Role, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, handlePostgresError(err))
	}

	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	const op = "storage.pg.UserRepository.FindByEmail"

	query, args, err := r.db.Builder.
		Select("id", "email", "password_hash", "role", "created_at").
		From("users").
		Where(squirrel.Eq{"email": email}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var user entity.User

	err = r.db.Pool.QueryRow(ctx, query, args...).Scan(
		&user.UserID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, handlePostgresError(err))
	}

	return &user, nil
}
