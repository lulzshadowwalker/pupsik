package handler

import (
	"log/slog"
	"net/http"
)

func MakeHandler(h func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			slog.ErrorContext(r.Context(), "Internal Server Error", "err", err)
			// TODO: render 500 page
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
