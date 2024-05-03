package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/josenymad/boulder-api/config"
	"github.com/josenymad/boulder-api/routes"
	_ "github.com/lib/pq"
)

var testFlag = flag.Bool("test", false, "Run in test mode")

func connectDB(test bool) {

	envVars := config.GetEnvVars(test)

	portInt, strconvErr := strconv.Atoi(envVars.Port)
	if strconvErr != nil {
		log.Fatal("Failed to convert port to integer:", strconvErr)
		return
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		envVars.Host, portInt, envVars.User, envVars.Password, envVars.Name)

	var err error
	config.DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = config.DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to the database")
}

func main() {
	flag.Parse()

	if *testFlag {
		fmt.Println("Running in test mode")
		connectDB(true)
	} else {
		fmt.Println("Running in normal mode")
		connectDB(false)
	}

	defer config.DB.Close()

	router := gin.Default()

	router.GET("/health", routes.HealthCheckHandler)

	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
