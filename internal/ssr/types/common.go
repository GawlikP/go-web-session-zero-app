package types

type StandardResponse struct {
	Success     bool   `json:"success"`
	RedirectURL string `json:"redirect_url"`
	Message     string `json:"message,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Fields  map[string]string
}
