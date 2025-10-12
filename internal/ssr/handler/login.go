package handler

import (
	"bytes"
	"log/slog"
	"net/http"
	"session-zero-app/web/templates"
)

func HandleLoginGet() http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := r.Context().Value("logger").(*slog.Logger)
		log.Info("[2] Starting login!")
		buf := new(bytes.Buffer)
		component := templates.LoginPage("")
		err := component.Render(r.Context(), buf)

		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			return
		}
		log.Info("[3] Got Valid Component to render!: LoginPage")
		w.Write(buf.Bytes())
	})
}
