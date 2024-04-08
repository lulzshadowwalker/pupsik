package handler

import (
	"net/http"

	"github.com/lulzshadowwalker/pupsik/view/home"
)

func HandleHomeIndex(w http.ResponseWriter, r *http.Request) error {
	return home.Index().Render(r.Context(), w)
}
