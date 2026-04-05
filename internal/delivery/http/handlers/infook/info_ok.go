package infook

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/api"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.Status(r, http.StatusOK)
		render.JSON(w, r, api.OK())
	}
}
