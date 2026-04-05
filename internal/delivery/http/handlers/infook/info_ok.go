package infook

import (
	"net/http"

	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/api"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		api.SendEmpty(r, http.StatusOK)
	}
}
