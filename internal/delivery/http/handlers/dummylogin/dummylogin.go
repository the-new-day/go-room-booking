package dummylogin

import (
	"log/slog"
	"net/http"

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
	api.ErrorResponse
	Token string `json:"token"`
}

func New(logger *slog.Logger, jwtManager *auth.JwtManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, ok := api.DecodeRequest[Request](logger, w, r)
		if !ok {
			return
		}

		userID := dummyUserID
		if req.Role == string(entity.RoleAdmin) {
			userID = dummyAdminID
		}

		accessToken, err := jwtManager.CreateToken(userID, req.Role)
		if err != nil {
			logger.Error("failed to generate token", sl.Err(err))

			api.SendInternalError(w, r, "failed to generate token")
			return
		}

		logger.Debug(
			"dummy login succeded",
			slog.String("role", req.Role),
		)

		api.SendOK(w, r, Response{
			Token: accessToken,
		})
	}
}
