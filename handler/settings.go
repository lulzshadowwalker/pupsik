package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/lulzshadowwalker/pupsik/database"
	"github.com/lulzshadowwalker/pupsik/types"
	"github.com/lulzshadowwalker/pupsik/utils"
	"github.com/lulzshadowwalker/pupsik/view/settings"
)

func HandleSettingsIndex(w http.ResponseWriter, r *http.Request) error {
	return render(w, r, settings.Index())
}

func HandleAccountSetupIndex(w http.ResponseWriter, r *http.Request) error {
	user := utils.GetUserFromContext(r.Context())
	if user == (types.User{}) {
		panic("account setup page should not be accessed by unauthenticated users")
	}

	_, err := database.GetAccountByUserID(r.Context(), user.ID)
	if err != nil && !errors.Is(err, qrm.ErrNoRows) {
		return fmt.Errorf("failed to get account by user id because %w", err)
	}

	return render(w, r, settings.AccountSetupIndex())
}

func HandleAccountSetupCreate(w http.ResponseWriter, r *http.Request) error {
	var errors settings.AccountSetupErrors
	params := settings.AccountSetupParams{
		Username: r.FormValue("username"),
	}

	if params.Username == "" {
		errors.Username = "Please provide a username"
	}

	if errors != (settings.AccountSetupErrors{}) {
		return render(w, r, settings.AccountSetupForm(params, errors))
	}

	user := utils.GetUserFromContext(r.Context())
	if user == (types.User{}) {
		panic("account setup page should not be accessed by unauthenticated users")
	}

	account := types.Account{
		UserID:   user.ID,
		Username: params.Username,
	}
	_, err := database.CreateAccount(r.Context(), account, nil)
	if err != nil {
		return err
	}

	return HxRedirect(w, r, "/settings")
}
