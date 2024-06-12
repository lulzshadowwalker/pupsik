package utils

import (
	"fmt"
	"net/url"

	"github.com/lulzshadowwalker/pupsik/config"
)

func GetURL(s string) (string, error) {
	res, err := url.JoinPath(config.GetAppURL(), s)
	if err != nil {
		return "", fmt.Errorf("failed to generate full url because %w", err)
	}

	return res, nil
}
