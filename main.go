package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/josenymad/boulder-api/config"
	_ "github.com/lib/pq"
)

func initDB() {
	envVars := config.GetEnvVars()

	portInt, strconvErr := strconv.Atoi(envVars.Port)
	if strconvErr != nil {
		log.Fatal("Failed to convert port to integer:", strconvErr)
		return
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		envVars.Host, portInt, envVars.User, envVars.Password, envVars.Name)

	var err error
	config.Db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = config.Db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to the database")
}

func main() {
	initDB()
	defer config.Db.Close()

	server := gin.Default()

	err := server.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}