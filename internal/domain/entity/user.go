package entity

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
	UserID       uuid.UUID `db:"user_id"`
	Email        string    `db:"email"`
	Role         UserRole  `db:"role"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
}

func IsValidRole(role UserRole) bool {
	switch role {
	case RoleAdmin, RoleUser:
		return true
	}

	return false
}
