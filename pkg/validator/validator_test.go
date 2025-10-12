package validator

import (
	"testing"
)

// TestNew verifies validator initialization
func TestNew(t *testing.T) {
	t.Parallel()

	v := New()

	if v == nil {
		t.Fatal("New() returned nil")
	}

	if v.Errors == nil {
		t.Error("Errors slice not initialized")
	}

	if len(v.Errors) != 0 {
		t.Errorf("Expected 0 errors, got %d", len(v.Errors))
	}

	if !v.Valid() {
		t.Error("New validator should be valid")
	}
}

// TestRequired tests the Required validation
func TestRequired(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		field     string
		value     string
		wantError bool
	}{
		{
			name:      "empty string",
			field:     "email",
			value:     "",
			wantError: true,
		},
		{
			name:      "whitespace only",
			field:     "name",
			value:     "   ",
			wantError: false, // Current implementation doesn't trim
		},
		{
			name:      "valid value",
			field:     "username",
			value:     "john",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // Each subtest runs in parallel

			v := New()
			v.Required(tt.field, tt.value)

			hasError := !v.Valid()
			if hasError != tt.wantError {
				t.Errorf("Required(%q, %q): hasError = %v, want %v",
					tt.field, tt.value, hasError, tt.wantError)
			}

			if tt.wantError && len(v.Errors) == 0 {
				t.Error("Expected error but got none")
			}

			if tt.wantError {
				if v.Errors[0].Field != tt.field {
					t.Errorf("Error field = %q, want %q", v.Errors[0].Field, tt.field)
				}
				if v.Errors[0].Message == "" {
					t.Error("Error message is empty")
				}
			}
		})
	}
}

// TestEmail tests email validation
func TestEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		email     string
		wantError bool
	}{
		{
			name:      "valid email",
			email:     "user@example.com",
			wantError: false,
		},
		{
			name:      "valid email with subdomain",
			email:     "user@mail.example.com",
			wantError: false,
		},
		{
			name:      "valid email with plus",
			email:     "user+tag@example.com",
			wantError: false,
		},
		{
			name:      "empty string",
			email:     "",
			wantError: false, // Email validation skips empty strings
		},
		{
			name:      "missing @",
			email:     "userexample.com",
			wantError: true,
		},
		{
			name:      "missing domain",
			email:     "user@",
			wantError: true,
		},
		{
			name:      "missing local part",
			email:     "@example.com",
			wantError: true,
		},
		{
			name:      "multiple @",
			email:     "user@@example.com",
			wantError: true,
		},
		{
			name:      "spaces",
			email:     "user @example.com",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := New()
			v.Email("email", tt.email)

			hasError := !v.Valid()
			if hasError != tt.wantError {
				t.Errorf("Email(%q): hasError = %v, want %v",
					tt.email, hasError, tt.wantError)
			}
		})
	}
}

// TestMinLength tests minimum length validation
func TestMinLength(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		value     string
		min       int
		wantError bool
	}{
		{
			name:      "empty string",
			value:     "",
			min:       5,
			wantError: false, // Skips empty
		},
		{
			name:      "exactly minimum",
			value:     "12345",
			min:       5,
			wantError: false,
		},
		{
			name:      "above minimum",
			value:     "123456",
			min:       5,
			wantError: false,
		},
		{
			name:      "below minimum",
			value:     "1234",
			min:       5,
			wantError: true,
		},
		{
			name:      "way below minimum",
			value:     "1",
			min:       10,
			wantError: true,
		},
		{
			name:      "unicode characters",
			value:     "你好世界",
			min:       3,
			wantError: false, // 4 characters
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := New()
			v.MinLength("field", tt.value, tt.min)

			hasError := !v.Valid()
			if hasError != tt.wantError {
				t.Errorf("MinLength(%q, %d): hasError = %v, want %v",
					tt.value, tt.min, hasError, tt.wantError)
			}

			if tt.wantError && len(v.Errors) > 0 {
				expectedMsg := "Must be at least"
				if v.Errors[0].Message[:len(expectedMsg)] != expectedMsg {
					t.Errorf("Error message = %q, want to start with %q",
						v.Errors[0].Message, expectedMsg)
				}
			}
		})
	}
}

// TestMaxLength tests maximum length validation
func TestMaxLength(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		value     string
		max       int
		wantError bool
	}{
		{
			name:      "empty string",
			value:     "",
			max:       5,
			wantError: false,
		},
		{
			name:      "exactly maximum",
			value:     "12345",
			max:       5,
			wantError: false,
		},
		{
			name:      "below maximum",
			value:     "1234",
			max:       5,
			wantError: false,
		},
		{
			name:      "above maximum",
			value:     "123456",
			max:       5,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := New()
			v.MaxLength("field", tt.value, tt.max)

			hasError := !v.Valid()
			if hasError != tt.wantError {
				t.Errorf("MaxLength(%q, %d): hasError = %v, want %v",
					tt.value, tt.max, hasError, tt.wantError)
			}
		})
	}
}

// TestMultipleValidations tests accumulating multiple errors
func TestMultipleValidations(t *testing.T) {
	t.Parallel()

	v := New()

	v.Required("email", "")
	v.Required("password", "")
	v.MinLength("username", "ab", 3)

	if len(v.Errors) != 3 {
		t.Errorf("Expected 3 errors, got %d", len(v.Errors))
	}

	if v.Valid() {
		t.Error("Validator should not be valid with errors")
	}
}

// TestFieldErrors tests the FieldErrors map generation
func TestFieldErrors(t *testing.T) {
	t.Parallel()

	v := New()

	v.AddError("email", "Invalid email")
	v.AddError("password", "Too short")
	v.AddError("username", "Required")

	fieldErrors := v.FieldErrors()

	if len(fieldErrors) != 3 {
		t.Errorf("Expected 3 field errors, got %d", len(fieldErrors))
	}

	if fieldErrors["email"] != "Invalid email" {
		t.Errorf("email error = %q, want %q", fieldErrors["email"], "Invalid email")
	}

	if fieldErrors["password"] != "Too short" {
		t.Errorf("password error = %q, want %q", fieldErrors["password"], "Too short")
	}

	if fieldErrors["username"] != "Required" {
		t.Errorf("username error = %q, want %q", fieldErrors["username"], "Required")
	}
}

// TestFieldErrorsFirstOnly tests that only first error per field is kept
func TestFieldErrorsFirstOnly(t *testing.T) {
	t.Parallel()

	v := New()

	v.AddError("password", "First error")
	v.AddError("password", "Second error")
	v.AddError("password", "Third error")

	fieldErrors := v.FieldErrors()

	if len(fieldErrors) != 1 {
		t.Errorf("Expected 1 field error, got %d", len(fieldErrors))
	}

	if fieldErrors["password"] != "First error" {
		t.Errorf("password error = %q, want %q", fieldErrors["password"], "First error")
	}
}

// TestCompleteValidationFlow tests a realistic validation scenario
func TestCompleteValidationFlow(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		email         string
		password      string
		wantValid     bool
		wantErrorsLen int
	}{
		{
			name:          "all valid",
			email:         "user@example.com",
			password:      "password123",
			wantValid:     true,
			wantErrorsLen: 0,
		},
		{
			name:          "empty email",
			email:         "",
			password:      "password123",
			wantValid:     false,
			wantErrorsLen: 1,
		},
		{
			name:          "invalid email",
			email:         "notanemail",
			password:      "password123",
			wantValid:     false,
			wantErrorsLen: 1,
		},
		{
			name:          "password too short",
			email:         "user@example.com",
			password:      "123",
			wantValid:     false,
			wantErrorsLen: 1,
		},
		{
			name:          "both invalid",
			email:         "",
			password:      "123",
			wantValid:     false,
			wantErrorsLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			v := New()

			v.Required("email", tt.email)
			v.Email("email", tt.email)
			v.Required("password", tt.password)
			v.MinLength("password", tt.password, 8)

			if v.Valid() != tt.wantValid {
				t.Errorf("Valid() = %v, want %v", v.Valid(), tt.wantValid)
			}

			if len(v.Errors) != tt.wantErrorsLen {
				t.Errorf("Got %d errors, want %d", len(v.Errors), tt.wantErrorsLen)
				for _, err := range v.Errors {
					t.Logf("  - %s: %s", err.Field, err.Message)
				}
			}
		})
	}
}

// BenchmarkValidator benchmarks validator performance
func BenchmarkValidator(b *testing.B) {
	b.Run("Required", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v := New()
			v.Required("field", "value")
		}
	})

	b.Run("Email", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v := New()
			v.Email("email", "user@example.com")
		}
	})

	b.Run("CompleteValidation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v := New()
			v.Required("email", "user@example.com")
			v.Email("email", "user@example.com")
			v.Required("password", "password123")
			v.MinLength("password", "password123", 8)
			_ = v.Valid()
			_ = v.FieldErrors()
		}
	})
}
