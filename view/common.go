package view

import (
	"context"

	"github.com/lulzshadowwalker/pupsik/types"
	"github.com/lulzshadowwalker/pupsik/utils"
)

func GetUser(ctx context.Context) types.User {
	u := utils.GetUserFromContext(ctx)
	return u
}
