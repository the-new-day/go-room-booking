package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/internships-backend/test-backend-the-new-day/pkg/logger/sl"
)

var (
	ErrFailedToDecodeRequest = errors.New("failed to decode request")
	ErrInvalidRequest        = errors.New("invalid request")
)

func DecodeRequest[Request any](logger *slog.Logger, w http.ResponseWriter, r *http.Request) (Request, bool) {
	var req Request
	err := render.DecodeJSON(r.Body, &req)

	if err != nil {
		logger.Error("failed to decode request body", sl.Err(err))

		SendBadRequest(w, r, "failed to decode request")
		return req, false
	}

	logger.Debug("request body decoded", slog.Any("request", req))

	if err := validator.New().Struct(req); err != nil {
		logger.Error("invalid request", sl.Err(err))

		SendBadRequest(w, r, ValidationError(err.(validator.ValidationErrors)))
		return req, false
	}

	return req, true
}
