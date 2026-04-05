package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
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

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}

func SendBadRequest(w http.ResponseWriter, r *http.Request, resp Response) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, resp)
}

func SendInternalServerError(w http.ResponseWriter, r *http.Request, resp Response) {
	render.Status(r, http.StatusInternalServerError)
	render.JSON(w, r, resp)
}

func SendNotFound(w http.ResponseWriter, r *http.Request, resp Response) {
	render.Status(r, http.StatusNotFound)
	render.JSON(w, r, resp)
}
