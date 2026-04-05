package booking

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/api"
	"github.com/internships-backend/test-backend-the-new-day/internal/domain/entity"
	"github.com/internships-backend/test-backend-the-new-day/pkg/logger/sl"
)

type PaginationResponse struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

type ListResponse struct {
	Bookings   []BookingResponse  `json:"bookings"`
	Pagination PaginationResponse `json:"pagination"`
}

type BookingLister interface {
	ListBookings(ctx context.Context, page, pageSize int) ([]*entity.Booking, int, error)
}

func NewListHandler(logger *slog.Logger, bookingLister BookingLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, pageSize, ok := parsePagination(r)
		if !ok {
			api.SendBadRequest(w, r, "invalid pagination")
			return
		}

		bookings, total, err := bookingLister.ListBookings(r.Context(), page, pageSize)
		if err != nil {
			logger.Error("failed to list bookings", sl.Err(err))
			api.SendInternalError(w, r, "failed to list bookings")
			return
		}

		api.SendOK(w, r, ListResponse{
			Bookings: mapBookingsToDto(bookings),
			Pagination: PaginationResponse{
				Page:     page,
				PageSize: pageSize,
				Total:    total,
			},
		})
	}
}

func parsePagination(r *http.Request) (int, int, bool) {
	page := 1
	pageSize := 20

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		val, err := strconv.Atoi(pageStr)
		if err != nil || val < 1 {
			return 0, 0, false
		}
		page = val
	}

	if pageSizeStr := r.URL.Query().Get("pageSize"); pageSizeStr != "" {
		val, err := strconv.Atoi(pageSizeStr)
		if err != nil || val < 1 || val > 100 {
			return 0, 0, false
		}
		pageSize = val
	}

	return page, pageSize, true
}
