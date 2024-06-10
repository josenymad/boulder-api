package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/josenymad/boulder-api/types"
)

var DB *sql.DB

const DefaultPort string = ":8080"

func ConnectDB(version string) error {
	envVars := GetEnvVars(version)

	var err error

	portInt, err := strconv.Atoi(envVars.Port)
	if err != nil {
		return fmt.Errorf("failed to convert port to integer: %v", err)
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		envVars.Host, portInt, envVars.User, envVars.Password, envVars.Name)

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	fmt.Println("Connected to the database")
	return nil
}

func GetEnvVars(version string) types.EnvVars {
	switch version {
	case "test":
		err := godotenv.Load("../.env.test")
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	case "development":
		err := godotenv.Load(".env.test")
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	default:
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}

	envVars := types.EnvVars{
		Host:     os.Getenv("HOST"),
		Port:     os.Getenv("PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}

	return envVars
}
