package view

import (
	"context"

	"github.com/lulzshadowwalker/pupsik/types"
	"github.com/lulzshadowwalker/pupsik/utils"
)

func GetUser(ctx context.Context) types.User {
	user, _ := utils.GetUserFromContext(ctx)
	return user
}
