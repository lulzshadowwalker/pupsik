package handler

import (
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gorilla/sessions"
	"github.com/lulzshadowwalker/pupsik/config"
	"github.com/lulzshadowwalker/pupsik/view/info"
)

type SessionKey string

const (
	SessionUserKey         = "user"
	SessionNotificationKey = "notifications"
	SessionToastKey        = "toast"
	SessionAccessTokenKey  = "access_token"
)

var Store = sessions.NewCookieStore([]byte(config.GetSessionKey()))

func Make(h func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			slog.ErrorContext(r.Context(), "Internal Server Error", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			render(w, r, info.InternalServerError())
		}
	}
}

func render(w http.ResponseWriter, r *http.Request, component templ.Component) error {
	return component.Render(r.Context(), w)
}

func HxRedirect(w http.ResponseWriter, r *http.Request, to string) error {
	if r.Header.Get("HX-Request") != "" {
		w.Header().Set("HX-Redirect", to)
		w.WriteHeader(http.StatusSeeOther)
		return nil
	}

	http.Redirect(w, r, to, http.StatusSeeOther)
	return nil
}
