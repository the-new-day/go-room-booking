package booking

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/internships-backend/test-backend-the-new-day/internal/domain"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain/entity"
	"github.com/internships-backend/test-backend-the-new-day/internal/storage"
)

type BookingRepository interface {
	Create(ctx context.Context, slotID, userID string, conferenceLink *string) (*entity.Booking, error)
	FindByID(ctx context.Context, bookingID string) (*entity.Booking, error)
	Cancel(ctx context.Context, bookingID string) (*entity.Booking, error)
	List(ctx context.Context, offset, limit int) ([]*entity.Booking, int, error)
	ListByUserFuture(ctx context.Context, userID string, now time.Time) ([]*entity.Booking, error)
}

type SlotRepository interface {
	GetByID(ctx context.Context, slotID string) (*entity.Slot, error)
}

type UseCase struct {
	bookingRepo BookingRepository
	slotRepo    SlotRepository
	now         func() time.Time
}

func New(bookingRepo BookingRepository, slotRepo SlotRepository) *UseCase {
	return &UseCase{
		bookingRepo: bookingRepo,
		slotRepo:    slotRepo,
		now:         func() time.Time { return time.Now().UTC() },
	}
}

func (u *UseCase) CreateBooking(ctx context.Context, slotID, userID string, createConferenceLink bool) (*entity.Booking, error) {
	const op = "usecase.booking.CreateBooking"

	slot, err := u.slotRepo.GetByID(ctx, slotID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w", op, domain.ErrSlotNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if slot.StartAt.Before(u.now()) {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrSlotInPast)
	}

	var conferenceLink *string
	_ = createConferenceLink

	booking, err := u.bookingRepo.Create(ctx, slotID, userID, conferenceLink)
	if err != nil {
		if errors.Is(err, storage.ErrSlotAlreadyBooked) {
			return nil, fmt.Errorf("%s: %w", op, domain.ErrSlotAlreadyBooked)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return booking, nil
}

func (u *UseCase) ListBookings(ctx context.Context, page, pageSize int) ([]*entity.Booking, int, error) {
	const op = "usecase.booking.ListBookings"

	offset := (page - 1) * pageSize

	bookings, total, err := u.bookingRepo.List(ctx, offset, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}

	return bookings, total, nil
}

func (u *UseCase) ListMyBookings(ctx context.Context, userID string) ([]*entity.Booking, error) {
	const op = "usecase.booking.ListMyBookings"

	bookings, err := u.bookingRepo.ListByUserFuture(ctx, userID, u.now())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return bookings, nil
}

func (u *UseCase) CancelBooking(ctx context.Context, bookingID, userID string) (*entity.Booking, error) {
	const op = "usecase.booking.CancelBooking"

	booking, err := u.bookingRepo.FindByID(ctx, bookingID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w", op, domain.ErrBookingNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if booking.UserID.String() != userID {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrForbidden)
	}

	if booking.Status == entity.BookingCancelled {
		return booking, nil
	}

	booking, err = u.bookingRepo.Cancel(ctx, bookingID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, fmt.Errorf("%s: %w", op, domain.ErrBookingNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return booking, nil
}
