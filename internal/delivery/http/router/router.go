package router

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/internships-backend/test-backend-the-new-day/internal/auth"
	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/handlers/dummylogin"
	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/handlers/infook"
	mw "github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/middleware"
)

type Router struct {
	*chi.Mux
}

func NewRouter(logger *slog.Logger, jwtManager *auth.JwtManager) *Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(mw.NewLoggerMiddleware(logger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Get("/_info", infook.New())
	r.Post("/dummyLogin", dummylogin.New(logger, jwtManager))

	return &Router{r}
}
