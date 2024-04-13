package types

import "github.com/nedpals/supabase-go"

type User struct {
	ID    string
	Email string
	Role  string
}

func NewUserFromSupabaseUser(user supabase.User) User {
	return User{
		ID:    user.ID,
		Role:  user.Role,
		Email: user.Email,
	}
}
