package router

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	authjwt "github.com/internships-backend/test-backend-the-new-day/internal/auth"
	authhandlers "github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/handlers/auth"
	bookinghandlers "github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/handlers/booking"
	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/handlers/dummylogin"
	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/handlers/infook"
	roomhandlers "github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/handlers/room"
	schedulehandlers "github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/handlers/schedule"
	slothandlers "github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/handlers/slot"
	mw "github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/middleware"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain/entity"
	"github.com/internships-backend/test-backend-the-new-day/internal/usecase/auth"
	"github.com/internships-backend/test-backend-the-new-day/internal/usecase/booking"
	"github.com/internships-backend/test-backend-the-new-day/internal/usecase/room"
	"github.com/internships-backend/test-backend-the-new-day/internal/usecase/schedule"
	"github.com/internships-backend/test-backend-the-new-day/internal/usecase/slot"
)

type Router struct {
	*chi.Mux
}

func NewRouter(
	logger *slog.Logger,
	jwtManager *authjwt.JwtManager,
	authUseCase *auth.UseCase,
	roomUseCase *room.UseCase,
	scheduleUseCase *schedule.UseCase,
	slotUseCase *slot.UseCase,
	bookingUseCase *booking.UseCase,
) *Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(mw.LoggerMiddleware(logger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Get("/_info", infook.New())
	r.Post("/dummyLogin", dummylogin.New(logger, jwtManager))

	r.Post("/login", authhandlers.NewLoginHandler(logger, authUseCase))
	r.Post("/register", authhandlers.NewRegisterHandler(logger, authUseCase))

	// Group requiring authorization (any role)
	r.Group(func(r chi.Router) {
		r.Use(mw.AuthMiddleware(logger, jwtManager))

		r.Get("/rooms/list", roomhandlers.NewListHandler(logger, roomUseCase))
		r.Get("/rooms/{roomId}/slots/list", slothandlers.NewListHandler(logger, slotUseCase))
	})

	// Admin only endpoints
	r.Group(func(r chi.Router) {
		r.Use(mw.AuthMiddleware(logger, jwtManager))
		r.Use(mw.RoleMiddleware(entity.RoleAdmin))

		r.Post("/rooms/create", roomhandlers.NewCreateHandler(logger, roomUseCase))
		r.Post("/rooms/{roomId}/schedule/create", schedulehandlers.NewCreateHandler(logger, scheduleUseCase))
		r.Get("/bookings/list", bookinghandlers.NewListHandler(logger, bookingUseCase))
	})

	// User only endpoints
	r.Group(func(r chi.Router) {
		r.Use(mw.AuthMiddleware(logger, jwtManager))
		r.Use(mw.RoleMiddleware(entity.RoleUser))

		r.Post("/bookings/create", bookinghandlers.NewCreateHandler(logger, bookingUseCase))
		r.Get("/bookings/my", bookinghandlers.NewMyHandler(logger, bookingUseCase))
		r.Post("/bookings/{bookingId}/cancel", bookinghandlers.NewCancelHandler(logger, bookingUseCase))
	})

	return &Router{r}
}
