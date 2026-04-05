package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/internships-backend/test-backend-the-new-day/internal/domain"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain/entity"
	"github.com/internships-backend/test-backend-the-new-day/internal/storage"
)

var (
	ErrFailedToHashPassword = errors.New("failed to hash password")
)

const DefaultRoleUponRegistration = entity.RoleUser

type UserRepository interface {
	Create(ctx context.Context, email, passwordHash string, role entity.UserRole) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
}

type PasswordHasher interface {
	DoesMatch(hash, password string) bool
	Hash(password string) (string, error)
}

type TokenManager interface {
	CreateToken(userID, role string) (string, error)
}

type UseCase struct {
	repo           UserRepository
	tokenManager   TokenManager
	passwordHasher PasswordHasher
}

func New(repo UserRepository, tokenManager TokenManager, passwordHasher PasswordHasher) *UseCase {
	return &UseCase{
		repo:           repo,
		tokenManager:   tokenManager,
		passwordHasher: passwordHasher,
	}
}

func (u *UseCase) Login(ctx context.Context, email, password string) (string, error) {
	const op = "usecase.auth.Login"

	user, err := u.repo.FindByEmail(ctx, email)

	if errors.Is(err, storage.ErrNotFound) {
		return "", fmt.Errorf("%s: %w", op, domain.ErrEmailNotFound)
	} else if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if !u.passwordHasher.DoesMatch(user.PasswordHash, password) {
		return "", fmt.Errorf("%s: %w", op, domain.ErrInvalidPassword)
	}

	token, err := u.tokenManager.CreateToken(user.UserID.String(), string(user.Role))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (u *UseCase) Register(ctx context.Context, email, password string) (*entity.User, error) {
	const op = "usecase.auth.Register"

	passwordHash, err := u.passwordHasher.Hash(password)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, ErrFailedToHashPassword)
	}

	user, err := u.repo.Create(ctx, email, passwordHash, DefaultRoleUponRegistration)

	if errors.Is(err, storage.ErrAlreadyExists) {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrUserWithEmailAlreadyExists)
	} else if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
