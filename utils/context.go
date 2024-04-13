package utils

import (
	"context"

	"github.com/lulzshadowwalker/pupsik/types"
)

type ContextKey string

const UserContextKey ContextKey = "user"

func GetUserFromContext(ctx context.Context) types.User {
	user, ok := ctx.Value(UserContextKey).(types.User)
	if !ok {
		return types.User{}
	}

	return user
}
