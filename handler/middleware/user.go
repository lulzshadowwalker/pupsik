package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/lulzshadowwalker/pupsik/database"
	"github.com/lulzshadowwalker/pupsik/handler"
	"github.com/lulzshadowwalker/pupsik/pkg/supa"
	"github.com/lulzshadowwalker/pupsik/types"
	"github.com/lulzshadowwalker/pupsik/utils"
)

func WithUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "public") {
			next.ServeHTTP(w, r)
			return
		}

		session, _ := handler.Store.Get(r, handler.SessionUserKey)
		accessToken := session.Values[handler.SessionAccessTokenKey]
		if accessToken == "" || accessToken == nil {
			next.ServeHTTP(w, r)
			return
		}

		supahUser, err := supa.Client.Auth.User(r.Context(), accessToken.(string))
		if err != nil {
			slog.Error("Failed to get user", "err", err)
			next.ServeHTTP(w, r)
			return
		}

		user := types.NewUserFromSupabaseUser(*supahUser)
		account, err := database.GetAccountByUserID(r.Context(), user.ID)
		if err != nil {
			if !errors.Is(err, qrm.ErrNoRows) {
				slog.ErrorContext(r.Context(), "failed to get account because %w", err)
			}
		}
		user.Account = account

		ctx := context.WithValue(r.Context(), utils.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
