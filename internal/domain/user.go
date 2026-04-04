package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

type User struct {
	UserID    uuid.UUID `db:"user_id"`
	Email     string    `db:"email"`
	Role      UserRole  `db:"role"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}
