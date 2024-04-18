package model

import (
	md "github.com/lulzshadowwalker/pupsik/database/.gen/postgres/public/model"
	"github.com/lulzshadowwalker/pupsik/types"
)

type Account struct {
	Account md.Account
}

func (a Account) ToEntity() types.Account {
	return types.Account{
		ID:       int(a.Account.ID),
		UserID:   a.Account.UserID,
		Username: a.Account.UserName,
	}
}
