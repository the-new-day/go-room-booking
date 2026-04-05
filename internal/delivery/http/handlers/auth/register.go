package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
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
	UserID string `json:"id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

type Registerer interface {
	Register(ctx context.Context, email, password, role string) (*entity.User, error)
}

func NewRegisterHandler(logger *slog.Logger, registerer Registerer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.NewRegisterHandler"

		logger = logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req RegisterRequest
		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			logger.Error("failed to decode request body", sl.Err(err))

			api.SendBadRequest(w, r, "failed to decode request")
			return
		}

		logger.Debug("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			logger.Error("invalid request", sl.Err(err))

			api.SendBadRequest(w, r, api.ValidationError(err.(validator.ValidationErrors)))
			return
		}

		user, err := registerer.Register(r.Context(), req.Email, req.Password, req.Role)

		if errors.Is(err, domain.ErrUserWithEmailAlreadyExists) {
			logger.Debug("user with provided email already exists", slog.String("email", req.Email))

			api.SendBadRequest(w, r, "user with provided email already exists")
			return
		} else if err != nil {
			logger.Error("failed to register user", sl.Err(err))

			api.SendInternalError(w, r, "failed to register user")
			return
		}

		api.SendCreated(w, r, RegisterResponse{
			UserID: user.UserID.String(),
			Email:  user.Email,
			Role:   string(user.Role),
		})
	}
}
