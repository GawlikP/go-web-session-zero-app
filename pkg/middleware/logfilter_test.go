package middleware

import (
	"testing"
)

func TestDefaultLogFilter(t *testing.T) {
	t.Parallel()

	filter := DefaultLogFilter()

	if filter == nil {
		t.Fatal("DefaultLogFilter() returned nil")
	}

	if len(filter.SkipPrefixes) == 0 {
		t.Error("Expected default SkipPrefixes to be set")
	}

	if len(filter.SkipExact) == 0 {
		t.Error("Expected default SkipExact to be set")
	}

	hasWellKnown := false
	for _, prefix := range filter.SkipPrefixes {
		if prefix == "/.well-known" {
			hasWellKnown = true
			break
		}
	}
	if !hasWellKnown {
		t.Error("Expected SkipPrefixes to contain '/.well-known'")
	}

	hasFavicon := false
	for _, exact := range filter.SkipExact {
		if exact == "/favicon.ico" {
			hasFavicon = true
			break
		}
	}
	if !hasFavicon {
		t.Error("Expected SkipExact to contain '/favicon.ico'")
	}
}

func TestLogFilter_ShouldSkip(t *testing.T) {
	t.Parallel()

	filter := DefaultLogFilter()

	tests := []struct {
		name   string
		path   string
		status int
		want   bool
	}{
		// Exact path matches
		{
			name:   "favicon exact match",
			path:   "/favicon.ico",
			status: 200,
			want:   true,
		},
		{
			name:   "robots.txt exact match",
			path:   "/robots.txt",
			status: 200,
			want:   true,
		},
		{
			name:   "sitemap.xml exact match",
			path:   "/sitemap.xml",
			status: 404,
			want:   true,
		},
		{
			name:   "ads.txt exact match",
			path:   "/ads.txt",
			status: 200,
			want:   true,
		},
		{
			name:   "security.txt exact match",
			path:   "/security.txt",
			status: 200,
			want:   true,
		},

		// Prefix matches
		{
			name:   "well-known prefix",
			path:   "/.well-known/security.txt",
			status: 200,
			want:   true,
		},
		{
			name:   "well-known chrome devtools",
			path:   "/.well-known/appspecific/com.chrome.devtools.json",
			status: 200,
			want:   true,
		},

		// Static files with success status
		{
			name:   "static CSS with 200",
			path:   "/static/css/bulma.min.css",
			status: 200,
			want:   true,
		},
		{
			name:   "static CSS with 304",
			path:   "/static/css/bulma.min.css",
			status: 304,
			want:   true,
		},
		{
			name:   "static JS with 200",
			path:   "/static/js/alpine.min.js",
			status: 200,
			want:   true,
		},
		{
			name:   "static image with 200",
			path:   "/static/images/logo.png",
			status: 200,
			want:   true,
		},

		// Static files with error status should NOT skip
		{
			name:   "static CSS with 404",
			path:   "/static/css/missing.css",
			status: 404,
			want:   false,
		},
		{
			name:   "static JS with 500",
			path:   "/static/js/app.js",
			status: 500,
			want:   false,
		},

		// Application endpoints should NOT skip
		{
			name:   "root path success",
			path:   "/",
			status: 200,
			want:   false,
		},
		{
			name:   "API endpoint success",
			path:   "/api/login",
			status: 200,
			want:   false,
		},
		{
			name:   "API endpoint error",
			path:   "/api/login",
			status: 400,
			want:   false,
		},
		{
			name:   "API endpoint server error",
			path:   "/api/login",
			status: 500,
			want:   false,
		},
		{
			name:   "dashboard page",
			path:   "/dashboard",
			status: 200,
			want:   false,
		},
		{
			name:   "workshop page",
			path:   "/workshop",
			status: 200,
			want:   false,
		},

		// Edge cases
		{
			name:   "empty path",
			path:   "",
			status: 200,
			want:   false,
		},
		{
			name:   "path with query params",
			path:   "/api/login?redirect=/dashboard",
			status: 200,
			want:   false,
		},
		{
			name:   "static with query params (cache busting)",
			path:   "/static/css/app.css?v=123",
			status: 200,
			want:   true,
		},
		{
			name:   "almost favicon",
			path:   "/favicon.ico.bak",
			status: 200,
			want:   false,
		},
		{
			name:   "almost well-known",
			path:   "/not-well-known/test",
			status: 200,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := filter.ShouldSkip(tt.path, tt.status)
			if got != tt.want {
				t.Errorf("ShouldSkip(%q, %d) = %v, want %v",
					tt.path, tt.status, got, tt.want)
			}
		})
	}
}

func TestLogFilter_ShouldSkip_CustomFilter(t *testing.T) {
	t.Parallel()

	// Test with custom filter configuration
	customFilter := &LogFilter{
		SkipPrefixes: []string{
			"/admin/metrics",
			"/health",
		},
		SkipExact: []string{
			"/ping",
			"/version",
		},
		SkipStatus: []int{304, 204},
	}

	tests := []struct {
		name   string
		path   string
		status int
		want   bool
	}{
		{
			name:   "custom exact match",
			path:   "/ping",
			status: 200,
			want:   true,
		},
		{
			name:   "custom prefix match",
			path:   "/admin/metrics/cpu",
			status: 200,
			want:   true,
		},
		{
			name:   "health prefix",
			path:   "/health/ready",
			status: 200,
			want:   true,
		},
		{
			name:   "not in custom list",
			path:   "/api/users",
			status: 200,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := customFilter.ShouldSkip(tt.path, tt.status)
			if got != tt.want {
				t.Errorf("ShouldSkip(%q, %d) = %v, want %v",
					tt.path, tt.status, got, tt.want)
			}
		})
	}
}

func TestLogFilter_ShouldSkip_EmptyFilter(t *testing.T) {
	t.Parallel()

	// Test with empty filter (nothing should be skipped)
	emptyFilter := &LogFilter{
		SkipPrefixes: []string{},
		SkipExact:    []string{},
		SkipStatus:   []int{},
	}

	tests := []struct {
		name   string
		path   string
		status int
	}{
		{"favicon", "/favicon.ico", 200},
		{"robots", "/robots.txt", 200},
		{"api", "/api/login", 200},
		{"static", "/static/css/app.css", 200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Static files should still be skipped (hardcoded in function)
			got := emptyFilter.ShouldSkip(tt.path, tt.status)

			if tt.path == "/static/css/app.css" {
				// Static files are hardcoded to skip
				if !got {
					t.Errorf("ShouldSkip(%q, %d) = false, want true (static files always skip)",
						tt.path, tt.status)
				}
			} else {
				// Everything else should NOT skip
				if got {
					t.Errorf("ShouldSkip(%q, %d) = true, want false (empty filter)",
						tt.path, tt.status)
				}
			}
		})
	}
}

// Benchmark the filter performance
func BenchmarkLogFilter_ShouldSkip(b *testing.B) {
	filter := DefaultLogFilter()

	b.Run("exact_match", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			filter.ShouldSkip("/favicon.ico", 200)
		}
	})

	b.Run("prefix_match", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			filter.ShouldSkip("/.well-known/security.txt", 200)
		}
	})

	b.Run("static_file", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			filter.ShouldSkip("/static/css/app.css", 200)
		}
	})

	b.Run("no_match", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			filter.ShouldSkip("/api/login", 200)
		}
	})
}
