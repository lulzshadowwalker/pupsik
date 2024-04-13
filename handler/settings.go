package handler

import (
	"net/http"

	"github.com/lulzshadowwalker/pupsik/view/settings"
)

func HandleSettingsIndex(w http.ResponseWriter, r *http.Request) error {
	return render(w, r, settings.Index())
}
