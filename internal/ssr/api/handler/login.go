package handler

import (
	"log/slog"
	"net/http"
	"session-zero-app/internal/ssr/types"
	"session-zero-app/pkg/response"
	"session-zero-app/pkg/validator"
)

func HandleLoginPost() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var req types.LoginRequest
		log := ctx.Value("logger").(*slog.Logger)
		log.Info("API login attempt")
		if err := response.ParseJSON(r, &req); err != nil {
			log.Warn("Invalid JSON", "error", err)
			response.WriteError(w, http.StatusBadRequest, "Invalid request form")
			return
		}
		v := validator.New()
		v.Required("email", req.Email)
		v.Email("email", req.Email)
		v.Required("password", req.Password)
		v.MinLength("password", req.Password, 8)

		if !v.Valid() {
			log.Warn("Validation failed", "failed", v.FieldErrors())
			response.WriteValidationError(w, v.FieldErrors())
			return
		}

		log.Info("Validation Passed!", "email", req.Email)
		response.WriteError(w, http.StatusUnauthorized, "Invalid email or password")
	})
}
