package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/lulzshadowwalker/pupsik/pkg/supa"
	"github.com/lulzshadowwalker/pupsik/utils"
	"github.com/lulzshadowwalker/pupsik/view/auth"
	"github.com/nedpals/supabase-go"
)

func HandleLoginIndex(w http.ResponseWriter, r *http.Request) error {
	return auth.LoginIndex().Render(r.Context(), w)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) error {
	credentials := supabase.UserCredentials{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	errors := auth.LoginErrors{}

	if ok, message := utils.ValidateEmail(credentials.Email); !ok {
		errors.Email = message
	}

	if ok, message := utils.ValidatePassword(credentials.Password); !ok {
		errors.Password = message
	}

	if (errors != auth.LoginErrors{}) {
		return render(w, r, auth.LoginForm(credentials, errors))
	}

	details, err := supa.Client.Auth.SignIn(r.Context(), credentials)
	if err != nil {
		return fmt.Errorf("failed to login user because %w", err)
	}

	setAuthCookie(w, details.AccessToken)
	slog.Info("User Login")
	return HxRedirect(w, r, "/")
}

func HandleLoginWithGoogle(w http.ResponseWriter, r *http.Request) error {
	details, err := supa.Client.Auth.SignInWithProvider(supabase.ProviderSignInOptions{
		Provider:   "google",
		RedirectTo: "http://localhost:3000/auth/callback",
	})
	if err != nil {
		return fmt.Errorf("failed to login with Google because %w", err)
	}

	HxRedirect(w, r, details.URL)
	return nil
}

func HandleRegisterIndex(w http.ResponseWriter, r *http.Request) error {
	return auth.RegisterIndex().Render(r.Context(), w)
}

func HandleRegister(w http.ResponseWriter, r *http.Request) error {
	params := auth.RegisterParams{
		Email:                r.FormValue("email"),
		Password:             r.FormValue("password"),
		PasswordConfirmation: r.FormValue("password_confirmation"),
	}

	errors := auth.RegisterErrors{}

	if ok, message := utils.ValidateEmail(params.Email); !ok {
		errors.Email = message
	}

	if ok, message := utils.ValidatePassword(params.Password); !ok {
		errors.Password = message
	}

	if params.Password != params.PasswordConfirmation {
		if params.PasswordConfirmation == "" {
			errors.PasswordConfirmation = "please enter your password one more time"
		} else {
			errors.PasswordConfirmation = "passwords do not match"
		}
	}

	if (errors != auth.RegisterErrors{}) {
		return render(w, r, auth.RegisterForm(params, errors))
	}

	credentials := supabase.UserCredentials{
		Email:    params.Email,
		Password: params.Password,
	}

	user, err := supa.Client.Auth.SignUp(r.Context(), credentials)
	if err != nil {
		return fmt.Errorf("failed to register new user because %w", err)
	}

	_ = user
	// setAuthCookie(w, user.Ac)

	slog.Info("User Registered")

	HxRedirect(w, r, "/")
	return nil
}

func HandleLogout(w http.ResponseWriter, r *http.Request) error {
	setAuthCookie(w, "")
	HxRedirect(w, r, "/auth/login")
	return nil
}

func setAuthCookie(w http.ResponseWriter, accessToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
	})
}

func HandleAuthCallback(w http.ResponseWriter, r *http.Request) error {
	accessToken := r.URL.Query().Get("access_token")
	if accessToken == "" {
		return render(w, r, auth.CallbackScript())
	}

	setAuthCookie(w, accessToken)
	HxRedirect(w, r, "/")
	return nil
}
