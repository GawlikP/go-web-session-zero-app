package types

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success     bool   `json:"success"`
	RedirectURL string `json:"redirect_url"`
	Message     string `json:"message,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Fields  map[string]string
}
