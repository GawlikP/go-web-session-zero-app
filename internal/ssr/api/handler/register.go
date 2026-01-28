package handler

import (
	"log/slog"
	"net/http"
	// "session-zero-app/pkg/logger"
	"encoding/json"
	"time"
	"session-zero-app/internal/ssr/types"
	"session-zero-app/pkg/response"
	"session-zero-app/pkg/validator"
	rbit "session-zero-app/pkg/messaging/rabbitmq"
	rbitTypes "session-zero-app/pkg/messaging/types"
)

func HandleRegisterPost(rpool *rbit.ConnectionPool) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var req types.RegisterRequest
		log := ctx.Value("logger").(*slog.Logger)
		log.Info("API request attempt!")
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
		v.MustMatch("password_confirmation", req.PasswordConfirmation, req.Password)

		if !v.Valid() {
			log.Warn("Validation failed", "failed", v.FieldErrors())
			response.WriteValidationError(w, v.FieldErrors())
			return
		}

		log.Info("Validation Passed!", "email", req.Email)
		log.Info("Validation Passed!", "message", "Sending...")

		command := rbitTypes.RegisterAttemptCommand{
			Email: req.Email,
			Password: req.Password,
			PasswordConfirm: req.PasswordConfirmation,
			IPAddress: r.Header.Get("X-Forwarded-For"),
			UserAgent: r.Header.Get("User-Agent"),
			RequestID: ctx.Value("request_id").(string),
			Timestamp: time.Now(),
		}

		mBody, err := json.Marshal(command)
		if err != nil {
			log.Warn("JSON marshaling for the command", "failed", err.Error())
			response.WriteError(w, http.StatusBadRequest, "Invalid request form")
			return
		}
		err = rbit.Publish(ctx, rpool, rbit.AuthCommandsExchange, rbit.AuthRegisterRoute, mBody)
		if err != nil {
			log.Warn("Rabbitmq publisher", "failed", err.Error())
			response.WriteError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		resp := types.StandardResponse{
			Success: true,
			RedirectURL: "",
			Message: "Registered Successfully!",
		}
		response.WriteJSON(w, http.StatusOK, resp)
	})
}
