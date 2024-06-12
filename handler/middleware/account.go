package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/lulzshadowwalker/pupsik/database"
	"github.com/lulzshadowwalker/pupsik/utils"
)

func WithAccount(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "public") {
			next.ServeHTTP(w, r)
			return
		}

		user, err := utils.GetUserFromContext(r.Context())
		if err != nil {
			slog.ErrorContext(r.Context(), "failed to get user account", "err", "WithAccount middleware should only be used when having an authenticated user")
			next.ServeHTTP(w, r)
			return
		}

		account, err := database.GetAccountByUserID(r.Context(), user.ID)
		if err != nil {
			if errors.Is(err, qrm.ErrNoRows) {
				http.Redirect(w, r, "/settings/account/setup", http.StatusSeeOther)
				return
			}

			slog.ErrorContext(r.Context(), "failed to get user account", "err", err)
			next.ServeHTTP(w, r)
			return
		}

		user.Account = account
		ctx := context.WithValue(r.Context(), utils.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
