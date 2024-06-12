package middleware

import (
	"net/http"
	"strings"

	"github.com/lulzshadowwalker/pupsik/handler"
	"github.com/lulzshadowwalker/pupsik/utils"
)

func WithAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "public") {
			next.ServeHTTP(w, r)
			return
		}

		_, err := utils.GetUserFromContext(r.Context())
		if err != nil {
			handler.HxRedirect(w, r, "/auth/login")
			return
		}

		if utils.ContainsAny(r.URL.Path, "login", "register") {
			handler.HxRedirect(w, r, "/")
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
