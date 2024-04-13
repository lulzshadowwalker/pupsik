package config

import (
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("failed to initialize config because " + err.Error())
	}
}

func GetSupabaseUrl() string {
	return os.Getenv("SUPABASE_URL")
}

func GetSupabaseSecret() string {
	return os.Getenv("SUPABASE_KEY")
}

func GetDatabasePassword() string {
	return os.Getenv("DB_PASSWORD")
}
