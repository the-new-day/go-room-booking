package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type ErrorCode string

const (
	ErrorCodeInvalidRequest    ErrorCode = "INVALID_REQUEST"
	ErrorCodeUnauthorized      ErrorCode = "UNAUTHORIZED"
	ErrorCodeNotFound          ErrorCode = "NOT_FOUND"
	ErrorCodeRoomNotFound      ErrorCode = "ROOM_NOT_FOUND"
	ErrorCodeSlotNotFound      ErrorCode = "SLOT_NOT_FOUND"
	ErrorCodeSlotAlreadyBooked ErrorCode = "SLOT_ALREADY_BOOKED"
	ErrorCodeBookingNotFound   ErrorCode = "BOOKING_NOT_FOUND"
	ErrorCodeForbidden         ErrorCode = "FORBIDDEN"
	ErrorCodeScheduleExists    ErrorCode = "SCHEDULE_EXISTS"
	ErrorCodeInternalError     ErrorCode = "INTERNAL_ERROR"
)

type ErrorResponse struct {
	Error ErrorDetail `json:"error,omitzero"`
}

type ErrorDetail struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func NewErrorResponse(code ErrorCode, message string) ErrorResponse {
	return ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
}

func ValidationError(errs validator.ValidationErrors) string {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid URL", err.Field()))
		case "oneof":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must contain value of: %s", err.Field(), err.Param()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return strings.Join(errMsgs, ", ")
}

func SendEmpty(r *http.Request, code int) {
	render.Status(r, code)
}

func SendBadRequest(w http.ResponseWriter, r *http.Request, message string) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, NewErrorResponse(ErrorCodeInvalidRequest, message))
}

func SendUnauthorized(w http.ResponseWriter, r *http.Request, message string) {
	render.Status(r, http.StatusUnauthorized)
	render.JSON(w, r, NewErrorResponse(ErrorCodeUnauthorized, message))
}

func SendForbidden(w http.ResponseWriter, r *http.Request, message string) {
	render.Status(r, http.StatusForbidden)
	render.JSON(w, r, NewErrorResponse(ErrorCodeForbidden, message))
}

func SendNotFound(w http.ResponseWriter, r *http.Request, message string) {
	render.Status(r, http.StatusNotFound)
	render.JSON(w, r, NewErrorResponse(ErrorCodeNotFound, message))
}

func SendRoomNotFound(w http.ResponseWriter, r *http.Request, message string) {
	render.Status(r, http.StatusNotFound)
	render.JSON(w, r, NewErrorResponse(ErrorCodeRoomNotFound, message))
}

func SendSlotNotFound(w http.ResponseWriter, r *http.Request, message string) {
	render.Status(r, http.StatusNotFound)
	render.JSON(w, r, NewErrorResponse(ErrorCodeSlotNotFound, message))
}

func SendBookingNotFound(w http.ResponseWriter, r *http.Request, message string) {
	render.Status(r, http.StatusNotFound)
	render.JSON(w, r, NewErrorResponse(ErrorCodeBookingNotFound, message))
}

func SendSlotAlreadyBooked(w http.ResponseWriter, r *http.Request, message string) {
	render.Status(r, http.StatusConflict)
	render.JSON(w, r, NewErrorResponse(ErrorCodeSlotAlreadyBooked, message))
}

func SendScheduleExists(w http.ResponseWriter, r *http.Request, message string) {
	render.Status(r, http.StatusConflict)
	render.JSON(w, r, NewErrorResponse(ErrorCodeScheduleExists, message))
}

func SendInternalError(w http.ResponseWriter, r *http.Request, message string) {
	render.Status(r, http.StatusInternalServerError)
	render.JSON(w, r, NewErrorResponse(ErrorCodeInternalError, message))
}

func SendOK(w http.ResponseWriter, r *http.Request, data interface{}) {
	render.Status(r, http.StatusOK)
	render.JSON(w, r, data)
}

func SendCreated(w http.ResponseWriter, r *http.Request, data interface{}) {
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, data)
}

func SendNoContent(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNoContent)
	w.WriteHeader(http.StatusNoContent)
}
