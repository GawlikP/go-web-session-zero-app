package validator

import (
	"fmt"
	"net/mail"
)

type ValidationError struct {
	Field   string
	Message string
}

type Validator struct {
	Errors []ValidationError
}

func New() *Validator {
	return &Validator{
		Errors: make([]ValidationError, 0),
	}
}

func (v *Validator) AddError(field, message string) {
	v.Errors = append(v.Errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) FieldErrors() map[string]string {
	fields := make(map[string]string)
	for _, err := range v.Errors {
		if _, exists := fields[err.Field]; !exists {
			fields[err.Field] = err.Message
		}
	}
	return fields
}

func (v *Validator) Required(field, value string) {
	if value == "" {
		v.AddError(field, "This field is required")
	}
}

func (v *Validator) Email(field, value string) {
	if value == "" {
		return
	}
	_, err := mail.ParseAddress(value)
	if err != nil {
		v.AddError(field, "Invalid email format")
	}
}

func (v *Validator) MinLength(field, value string, m int) {
	if value == "" {
		return
	}
	if len(value) < m {
		v.AddError(field, fmt.Sprintf("Must be at least %d characters", m))
	}
}

func (v *Validator) MaxLength(field, value string, m int) {
	if len(value) > m {
		v.AddError(field, fmt.Sprintf("Must be at most %d characters", m))
	}
}

func (v *Validator) MustMatch(field string, value string, other string) {
		if value != other {
				v.AddError(field, "Values do not match")
		}
}
