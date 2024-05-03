package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/josenymad/boulder-api/types"
)

func GetEnvVars(test bool) types.EnvVars {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if test {
		return types.EnvVars{
			Host:     os.Getenv("HOST"),
			Port:     os.Getenv("PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Name:     os.Getenv("TEST_DB_NAME"),
		}
	}

	return types.EnvVars{
		Host:     os.Getenv("HOST"),
		Port:     os.Getenv("PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}
}
