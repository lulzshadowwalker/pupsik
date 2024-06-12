package handler

import (
	"net/http"

	"github.com/lulzshadowwalker/pupsik/database"
	"github.com/lulzshadowwalker/pupsik/utils"
	"github.com/lulzshadowwalker/pupsik/view/home"
)

func HandleHomeIndex(w http.ResponseWriter, r *http.Request) error {
	user, err := utils.GetUserFromContext(r.Context())
	if err != nil {
		return err
	}

	images, err := database.GetImagesByUserID(r.Context(), user.ID, nil)
	if err != nil {
		return err
	}

	return render(w, r, home.Index(home.IndexParams{Images: images}))
}
