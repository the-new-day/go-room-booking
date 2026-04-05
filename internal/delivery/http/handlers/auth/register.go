package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/api"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain/entity"
	"github.com/internships-backend/test-backend-the-new-day/pkg/logger/sl"
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role" validate:"required,oneof=admin user"`
}

type RegisterResponse struct {
	User UserResponse `json:"user"`
}

type UserResponse struct {
	UserID    string  `json:"id"`
	Email     string  `json:"email"`
	Role      string  `json:"role"`
	CreatedAt *string `json:"createdAt"`
}

type Registerer interface {
	Register(ctx context.Context, email, password, role string) (*entity.User, error)
}

func NewRegisterHandler(logger *slog.Logger, registerer Registerer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, ok := api.DecodeRequest[RegisterRequest](logger, w, r)
		if !ok {
			return
		}

		user, err := registerer.Register(r.Context(), req.Email, req.Password, req.Role)

		if errors.Is(err, domain.ErrUserWithEmailAlreadyExists) {
			logger.Debug("user with provided email already exists", slog.String("email", req.Email))

			api.SendBadRequest(w, r, "user with provided email already exists")
			return
		} else if errors.Is(err, domain.ErrInvalidRole) {
			api.SendBadRequest(w, r, "invalid role")
			return
		} else if err != nil {
			logger.Error("failed to register user", sl.Err(err))

			api.SendInternalError(w, r, "failed to register user")
			return
		}

		api.SendCreated(w, r, RegisterResponse{
			User: UserResponse{
				UserID:    user.UserID.String(),
				Email:     user.Email,
				Role:      string(user.Role),
				CreatedAt: nil,
			},
		})
	}
}
