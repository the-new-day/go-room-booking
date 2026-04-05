package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/internships-backend/test-backend-the-new-day/config"
	authjwt "github.com/internships-backend/test-backend-the-new-day/internal/auth"
	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/router"
	"github.com/internships-backend/test-backend-the-new-day/internal/storage/pg"
	"github.com/internships-backend/test-backend-the-new-day/internal/usecase/auth"
	"github.com/internships-backend/test-backend-the-new-day/internal/usecase/booking"
	"github.com/internships-backend/test-backend-the-new-day/internal/usecase/room"
	"github.com/internships-backend/test-backend-the-new-day/internal/usecase/schedule"
	"github.com/internships-backend/test-backend-the-new-day/internal/usecase/slot"
	"github.com/internships-backend/test-backend-the-new-day/pkg/hasher"
	"github.com/internships-backend/test-backend-the-new-day/pkg/logger/sl"
	"github.com/internships-backend/test-backend-the-new-day/pkg/postgres"
)

const (
	logLevelDebug   = "debug"
	logLevelInfo    = "info"
	logLevelDiscard = "discard"
)

func Run(cfg *config.Config) {
	// Logger
	logger := SetupLogger(cfg.Log.Level)
	logger.Info("starting server")
	logger.Debug("debug messages are enabled")

	// Graceful shutdown signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Infrastructure
	logger.Info("connecting to Postgres")
	db := SetupDatabase(cfg.Postgres.DSN(), cfg.Postgres.MaxPoolSize)
	jwtManager := authjwt.NewJwtManager(cfg.JwtConfig.SignKey, cfg.JwtConfig.AccessTTL)
	passwordHasher := hasher.NewBcryptHasher()

	// Repositorires
	userRepo := pg.NewUserRepository(db)
	roomRepo := pg.NewRoomRepository(db)
	scheduleRepo := pg.NewScheduleRepository(db)
	slotRepo := pg.NewSlotRepository(db)
	bookingRepo := pg.NewBookingRepository(db)

	// Use cases
	authUseCase := auth.New(userRepo, jwtManager, passwordHasher)
	roomUseCase := room.New(roomRepo)
	scheduleUseCase := schedule.New(scheduleRepo, roomRepo, slotRepo)
	slotUseCase := slot.New(roomRepo, scheduleRepo, slotRepo)
	bookingUseCase := booking.New(bookingRepo, slotRepo)

	// HTTP server
	router := router.NewRouter(
		logger,
		jwtManager,
		authUseCase,
		roomUseCase,
		scheduleUseCase,
		slotUseCase,
		bookingUseCase,
	)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.HttpServer.Port),
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("failed to start server", sl.Err(err))
		}
	}()

	<-interrupt
	logger.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("failed to stop server", sl.Err(err))
	} else {
		logger.Info("server stopped")
	}

	logger.Info("closing Postgres")
	db.Close()
	logger.Info("Postgres closed")
}

func SetupLogger(level string) *slog.Logger {
	var logger *slog.Logger

	switch level {
	case logLevelDebug:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case logLevelInfo:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case logLevelDiscard:
		logger = sl.NewDiscardLogger()
	default:
		log.Fatalf("logger setup failed: unsupported level %q", level)
	}

	return logger
}

func SetupDatabase(dsn string, maxPoolSize int) *postgres.Postgres {
	pg, err := postgres.New(dsn, postgres.MaxPoolSize(maxPoolSize))
	if err != nil {
		log.Fatalf("database setup failed: %v", err)
	}
	return pg
}

func LoadConfig() *config.Config {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	return cfg
}
