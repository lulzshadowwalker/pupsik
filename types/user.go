package types

import (
	"github.com/google/uuid"
	"github.com/nedpals/supabase-go"
)

type User struct {
	ID      uuid.UUID
	Email   string
	Role    string
	Account Account
}

func NewUserFromSupabaseUser(user supabase.User) User {
	return User{
		ID:    uuid.MustParse(user.ID),
		Role:  user.Role,
		Email: user.Email,
	}
}
