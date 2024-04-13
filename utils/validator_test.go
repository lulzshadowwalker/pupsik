package utils

import "testing"

func TestValidateEmail(t *testing.T) {
	cases := []struct {
		Email string
		Valid bool
	}{
		{
			Email: "not an email",
			Valid: false,
		},
		{
			Email: "email@example.com",
			Valid: true,
		},
	}

	for _, kase := range cases {
		if ok, _ := ValidateEmail(kase.Email); ok != kase.Valid {
			t.Errorf("email %q, want %t got %t", kase.Email, ok, kase.Valid)
		}
	}
}
