package types

import "time"

type LoginAttemptCommand struct {
	Email string `json:"email"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	RequestID string `json:"request_id"`
	Timestamp time.Time `json:"timestamp"`
}

type RegisterAttemptCommand struct {
	Email string `json:"email"`
	Password string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	RequestID string `json:"request_id"`
	Timestamp time.Time `json:"timestamp"`
}
