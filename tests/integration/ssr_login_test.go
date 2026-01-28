package integration

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestLoginPageSSR(t *testing.T) {
	TestSetup(t)
	server := SetupTestServer(t)

	tests := []struct {
		name           string
		path           string
		wantStatus     int
		wantContains   []string
		wantNotContain []string
	}{
		{
			name:       "login page renders",
			path:       "/",
			wantStatus: 200,
			wantContains: []string{
				"<title>Login",
				"TTRPG Platform",
				"<form",
				"type=\"email\"",
				"type=\"password\"",
				"x-data=\"loginForm()\"",
			},
		},
		{
			name:       "includes Alpine.js",
			path:       "/",
			wantStatus: 200,
			wantContains: []string{
				"alpine.min.js",
			},
		},
		{
			name:       "includes Bulma CSS",
			path:       "/",
			wantStatus: 200,
			wantContains: []string{
				"bulma.min.css",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Get(server.URL + tt.path)
			if err != nil {
				t.Fatalf("GET request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("Status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}

			// Read body
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read body: %v", err)
			}
			body := string(bodyBytes)

			// Check contains
			for _, want := range tt.wantContains {
				if !strings.Contains(body, want) {
					t.Errorf("Body missing %q", want)
				}
			}

			// Check not contains
			for _, notWant := range tt.wantNotContain {
				if strings.Contains(body, notWant) {
					t.Errorf("Body should not contain %q", notWant)
				}
			}

			// Check Content-Type
			contentType := resp.Header.Get("Content-Type")
			if !strings.Contains(contentType, "text/html") {
				t.Errorf("Content-Type = %q, want text/html", contentType)
			}
		})
	}
}
