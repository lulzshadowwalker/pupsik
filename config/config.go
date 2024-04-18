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

func GetDatabaseHost() string {
	return os.Getenv("DATABASE_HOST")
}

func GetDatabaseUser() string {
	return os.Getenv("DATABASE_USER")
}

func GetDatabaseName() string {
	return os.Getenv("DATABASE_NAME")
}

func GetDatabasePort() string {
	return os.Getenv("DATABASE_PORT")
}

func GetDatabasePassword() string {
	return os.Getenv("DATABASE_PASSWORD")
}
