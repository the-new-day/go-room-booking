package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type tokenResponse struct {
	Token string `json:"token"`
}

type roomResponse struct {
	Room struct {
		ID          string  `json:"id"`
		Name        string  `json:"name"`
		Description *string `json:"description"`
		Capacity    *int    `json:"capacity"`
		CreatedAt   *string `json:"createdAt"`
	} `json:"room"`
}

type scheduleResponse struct {
	Schedule struct {
		ID         string `json:"id"`
		RoomID     string `json:"roomId"`
		DaysOfWeek []int  `json:"daysOfWeek"`
		StartTime  string `json:"startTime"`
		EndTime    string `json:"endTime"`
	} `json:"schedule"`
}

type slotsResponse struct {
	Slots []struct {
		ID     string `json:"id"`
		RoomID string `json:"roomId"`
		Start  string `json:"start"`
		End    string `json:"end"`
	} `json:"slots"`
}

type bookingResponse struct {
	Booking struct {
		ID             string  `json:"id"`
		SlotID         string  `json:"slotId"`
		UserID         string  `json:"userId"`
		Status         string  `json:"status"`
		ConferenceLink *string `json:"conferenceLink"`
		CreatedAt      *string `json:"createdAt"`
	} `json:"booking"`
}

func baseURL() string {
	if v := os.Getenv("E2E_BASE_URL"); v != "" {
		return v
	}
	return "http://localhost:8080"
}

func isServerUp(baseURL string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/_info", nil)
	if err != nil {
		return false
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func dummyLogin(t *testing.T, baseURL, role string) string {
	payload := map[string]any{"role": role}
	var res tokenResponse
	doJSON(t, http.MethodPost, baseURL+"/dummyLogin", "", payload, &res)
	require.NotEmpty(t, res.Token)
	return res.Token
}

func createRoom(t *testing.T, baseURL, token string) roomResponse {
	payload := map[string]any{
		"name": fmt.Sprintf("room-%s", uuid.New().String()),
	}
	var res roomResponse
	doJSON(t, http.MethodPost, baseURL+"/rooms/create", token, payload, &res)
	require.NotEmpty(t, res.Room.ID)
	return res
}

func createSchedule(t *testing.T, baseURL, token, roomID string) scheduleResponse {
	payload := map[string]any{
		"id":         uuid.New().String(),
		"roomId":     roomID,
		"daysOfWeek": []int{1, 2, 3, 4, 5, 6, 7},
		"startTime":  "09:00",
		"endTime":    "11:00",
	}
	var res scheduleResponse
	doJSON(t, http.MethodPost, fmt.Sprintf("%s/rooms/%s/schedule/create", baseURL, roomID), token, payload, &res)
	require.NotEmpty(t, res.Schedule.ID)
	return res
}

func listSlots(t *testing.T, baseURL, token, roomID string) slotsResponse {
	date := time.Now().UTC().Add(24 * time.Hour).Format("2006-01-02")
	var res slotsResponse
	doJSON(t, http.MethodGet, fmt.Sprintf("%s/rooms/%s/slots/list?date=%s", baseURL, roomID, date), token, nil, &res)
	return res
}

func createBooking(t *testing.T, baseURL, token, slotID string) bookingResponse {
	payload := map[string]any{"slotId": slotID}
	var res bookingResponse
	doJSON(t, http.MethodPost, baseURL+"/bookings/create", token, payload, &res)
	require.NotEmpty(t, res.Booking.ID)
	return res
}

func cancelBooking(t *testing.T, baseURL, token, bookingID string) bookingResponse {
	var res bookingResponse
	doJSON(t, http.MethodPost, fmt.Sprintf("%s/bookings/%s/cancel", baseURL, bookingID), token, nil, &res)
	return res
}

func doJSON(t *testing.T, method, url, token string, payload any, out any) {
	t.Helper()

	var body []byte
	if payload != nil {
		var err error
		body, err = json.Marshal(payload)
		require.NoError(t, err)
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	require.NoError(t, err)

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.True(t, resp.StatusCode >= 200 && resp.StatusCode < 300, "unexpected status: %d", resp.StatusCode)

	if out != nil {
		err = json.NewDecoder(resp.Body).Decode(out)
		require.NoError(t, err)
	}
}
