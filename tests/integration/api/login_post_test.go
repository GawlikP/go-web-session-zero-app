package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"session-zero-app/tests/integration"
	"testing"
)

func TestSSRLoginAPI(t *testing.T) {
	integration.TestSetup(t)
	server := integration.SetupTestServer(t)
	defer server.Close()

	tests := []struct {
		name       string
		payload    map[string]string
		wantStatus int
	}{
		{
			name: "valid credentials",
			payload: map[string]string{
				"email":    "user@example.com",
				"password": "password123",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "missing email",
			payload: map[string]string{
				"password": "password123",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid password",
			payload: map[string]string{
				"email":    "test@mail.com",
				"password": "1234567",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)

			resp, err := http.Post(
				server.URL+"/api/login",
				"application/json",
				bytes.NewReader(body),
			)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("Status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}
