package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/api"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain"
	"github.com/internships-backend/test-backend-the-new-day/pkg/logger/sl"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type Loginner interface {
	Login(ctx context.Context, email, password string) (string, error)
}

func NewLoginHandler(logger *slog.Logger, loginner Loginner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, ok := api.DecodeRequest[LoginRequest](logger, w, r)
		if !ok {
			return
		}

		token, err := loginner.Login(r.Context(), req.Email, req.Password)

		if errors.Is(err, domain.ErrEmailNotFound) {
			logger.Debug("email not found")

			api.SendUnauthorized(w, r, "email not found")
			return
		} else if errors.Is(err, domain.ErrInvalidPassword) {
			logger.Debug("invalid password")

			api.SendUnauthorized(w, r, "invalid password")
			return
		} else if err != nil {
			logger.Error("login failed", sl.Err(err))

			api.SendInternalError(w, r, "login failed")
			return
		}

		api.SendOK(w, r, LoginResponse{
			Token: token,
		})
	}
}
