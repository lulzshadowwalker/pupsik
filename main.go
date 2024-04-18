package main

import (
	"embed"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lulzshadowwalker/pupsik/handler"
	mw "github.com/lulzshadowwalker/pupsik/handler/middleware"
)

//go:embed public
var FS embed.FS

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(mw.WithUser)

	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.FS(FS))))

	r.Get("/auth/login", handler.Make(handler.HandleLoginIndex))
	r.Post("/auth/login", handler.Make(handler.HandleLogin))
	r.Get("/auth/register", handler.Make(handler.HandleRegisterIndex))
	r.Post("/auth/register", handler.Make(handler.HandleRegister))
	r.Post("/auth/logout", handler.Make(handler.HandleLogout))
	r.Get("/auth/callback", handler.Make(handler.HandleAuthCallback))
	r.Post("/auth/providers/google", handler.Make(handler.HandleLoginWithGoogle))

	r.Group(func(auth chi.Router) {
		auth.Use(mw.WithAuth)
		auth.Get("/settings/account/setup", handler.Make(handler.HandleAccountSetupIndex))
		auth.Post("/settings/account/setup", handler.Make(handler.HandleAccountSetupCreate))

		auth.Group(func(account chi.Router) {
			account.Use(mw.WithAccount)
			account.Get("/", handler.Make(handler.HandleHomeIndex))
			account.Get("/settings", handler.Make(handler.HandleSettingsIndex))
		})
	})

	slog.Info("server started", "port", ":3000")
	err := http.ListenAndServe(":3000", r)
	slog.Error("server shutdown", "err", err)
}
