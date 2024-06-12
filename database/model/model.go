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

type DBImage struct {
	Image md.Image
}

func (i DBImage) ToEntity() types.Image {
	img := types.Image{
		ID:        int(i.Image.ID),
		UserID:    i.Image.UserID,
		CreatedAt: i.Image.CreatedAt,
		Prompt:    i.Image.Prompt,
		Status:    types.ImageStatus(i.Image.Status),
		// DeletedAt: i.Image.DeletedAt,
	}

	if i.Image.URL != nil {
		img.URL = *i.Image.URL
	}

	return img
}
