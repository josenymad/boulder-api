package config

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var DB *sql.DB

type EnvVars struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func GetEnvVars(test bool) EnvVars {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if test {
		return EnvVars{
			Host:     os.Getenv("HOST"),
			Port:     os.Getenv("PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Name:     os.Getenv("TEST_DB_NAME"),
		}
	}

	return EnvVars{
		Host:     os.Getenv("HOST"),
		Port:     os.Getenv("PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}
}
