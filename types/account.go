package types

import "github.com/google/uuid"

type Account struct {
	ID       int
	UserID   uuid.UUID
	Username string
}
