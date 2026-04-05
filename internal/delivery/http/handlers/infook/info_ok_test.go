package infook

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/internships-backend/test-backend-the-new-day/internal/delivery/http/api"
	"github.com/stretchr/testify/require"
)

func TestInfoOkHandler(t *testing.T) {
	tests := []struct {
		name  string
		query map[string]string
	}{
		{
			name: "empty request",
		},
		{
			name: "request with query params",
			query: map[string]string{
				"param1": "value1",
				"param2": "value2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			endpoint := "/_info"
			if tt.query != nil {
				params := url.Values{}
				for name, val := range tt.query {
					params.Add(name, val)
				}
				endpoint += "?" + params.Encode()
			}

			rr := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, endpoint, nil)

			handler := New()
			handler.ServeHTTP(rr, r)

			var resp api.Response
			require.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))

			require.Equal(t, http.StatusOK, rr.Code)
			require.Equal(t, resp.Status, api.StatusOK)
		})
	}
}
