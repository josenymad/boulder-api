package config

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var Db *sql.DB

type EnvVars struct {
	Host string;
	Port string;
	User string;
	Password string;
	Name string;
}

func GetEnvVars () EnvVars {
	err := godotenv.Load()
	if err != nil {
        log.Fatal("Error loading .env file")
    }
	
	return EnvVars{
		Host: os.Getenv("HOST"),
		Port: os.Getenv("PORT"),
		User: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Name: os.Getenv("DB_NAME"),
	}
}