package main

import (
	"embed"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lulzshadowwalker/pupsik/handler"
)

//go:embed public
var FS embed.FS

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.FS(FS))))
	r.Get("/", handler.MakeHandler(handler.HandleHomeIndex))

	slog.Info("server started", "port", ":3000")
	err := http.ListenAndServe(":3000", r)
	slog.Error("server shutdown", "err", err)
}
