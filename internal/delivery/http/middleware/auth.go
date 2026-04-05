package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/internships-backend/test-backend-the-new-day/internal/auth"
	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/api"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	RoleKey   contextKey = "role"
)

func AuthMiddleware(logger *slog.Logger, jwtManager *auth.JwtManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				api.SendUnauthorized(w, r, api.Error("no Authorization header present"))
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" || strings.TrimSpace(parts[1]) == "" {
				api.SendUnauthorized(w, r, api.Error("invalid Authorization header"))
				return
			}

			claims, err := jwtManager.ParseToken(parts[1])
			if err != nil {
				api.SendUnauthorized(w, r, api.Error("invalid token"))
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, RoleKey, claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
