package utils

import (
	"fmt"
	"net/mail"
)

func ValidatePassword(password string) (ok bool, message string) {
	if password == "" {
		return false, "password cannot be empty"
	}

	const min = 8
	if len(password) < min {
		return false, fmt.Sprintf("password cannot be less than %d characters", min)
	}

	const max = 36
	if len(password) > max {
		return false, fmt.Sprintf("password cannot be longer than %d characters", max)
	}

	return true, ""
}

func ValidateEmail(email string) (ok bool, message string) {
	if _, err := mail.ParseAddress(email); err != nil {
		return false, "please enter a valid email address"
	}

	return true, ""
}
