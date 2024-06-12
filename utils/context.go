package utils

import (
	"context"
	"errors"

	"github.com/lulzshadowwalker/pupsik/types"
)

type ContextKey string

const (
	UserContextKey                 ContextKey = "user"
	SessionNotificationsContextKey ContextKey = "notifications"
)

func GetUserFromContext(ctx context.Context) (types.User, error) {
	user, ok := ctx.Value(UserContextKey).(types.User)
	if !ok {
		return types.User{}, errors.New("user key not found in context")
	}

	return user, nil
}
