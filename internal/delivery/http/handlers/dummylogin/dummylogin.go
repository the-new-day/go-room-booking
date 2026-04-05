package dummylogin

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/internships-backend/test-backend-the-new-day/internal/auth"
	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/api"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain/entity"
	"github.com/internships-backend/test-backend-the-new-day/pkg/logger/sl"
)

const (
	dummyAdminID = "00000000-0000-0000-0000-000000000001"
	dummyUserID  = "00000000-0000-0000-0000-000000000002"
)

type Request struct {
	Role string `json:"role" validate:"required,oneof=admin user"`
}

type Response struct {
	api.Response
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

func New(logger *slog.Logger, jwtManager *auth.JwtManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.dummylogin.New"

		logger = logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			logger.Error("failed to decode request body", sl.Err(err))

			api.SendBadRequest(w, r, api.Error("failed to decode request"))
			return
		}

		logger.Debug("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			logger.Error("invalid request", sl.Err(err))

			api.SendBadRequest(w, r, api.ValidationError(err.(validator.ValidationErrors)))
			return
		}

		logger.Info(
			"dummy login succeded",
			slog.String("role", req.Role),
		)

		userID := dummyUserID
		if req.Role == string(entity.RoleAdmin) {
			userID = dummyAdminID
		}

		accessToken, err := jwtManager.CreateToken(userID, req.Role)
		if err != nil {
			logger.Error("failed to generate token", sl.Err(err))

			api.SendInternalServerError(w, r, api.Error("failed to generate token"))
			return
		}

		render.JSON(w, r, Response{
			Response:    api.OK(),
			AccessToken: accessToken,
			TokenType:   "Bearer",
			ExpiresIn:   int64(jwtManager.AccessTokenTTL().Seconds()),
		})
	}
}
