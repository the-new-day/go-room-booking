package e2e

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestE2E_CreateRoomScheduleBooking(t *testing.T) {
	baseURL := baseURL()
	if !isServerUp(baseURL) {
		t.Skip("server is not reachable on " + baseURL)
	}

	adminToken := dummyLogin(t, baseURL, "admin")
	userToken := dummyLogin(t, baseURL, "user")

	room := createRoom(t, baseURL, adminToken)
	createSchedule(t, baseURL, adminToken, room.Room.ID)

	slots := listSlots(t, baseURL, adminToken, room.Room.ID)
	require.NotEmpty(t, slots.Slots)

	booking := createBooking(t, baseURL, userToken, slots.Slots[0].ID)
	require.Equal(t, slots.Slots[0].ID, booking.Booking.SlotID)
	require.Equal(t, "active", booking.Booking.Status)
}
