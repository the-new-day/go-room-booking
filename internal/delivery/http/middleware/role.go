package middleware

import (
	"net/http"

	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/api"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain/entity"
)

func RoleMiddleware(role entity.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			currentRole, ok := r.Context().Value(RoleKey).(string)

			if !ok || currentRole == "" {
				api.SendUnauthorized(w, r, "unauthorized")
				return
			}

			if currentRole != string(role) {
				api.SendForbidden(w, r, "no permission")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
