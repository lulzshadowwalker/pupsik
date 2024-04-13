package supa

import (
	"github.com/lulzshadowwalker/pupsik/config"
	"github.com/nedpals/supabase-go"
)

var Client *supabase.Client

func init() {
	url := config.GetSupabaseUrl()
	key := config.GetSupabaseSecret()

	Client = supabase.CreateClient(url, key)
}
