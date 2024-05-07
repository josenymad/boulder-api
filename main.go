package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

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
		log.Fatal("Failed to connect to the database:", err)
	}

	err = config.DB.Ping()
	if err != nil {
		log.Fatal("Failed to ping the database:", err)
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
	router.POST("/categories", routes.CreateCompetitionCategory)
	router.POST("/rounds", routes.CreateRound)
	router.POST("/competitors", routes.CreateCompetitor)

	// Graceful shutdown
	srv := &http.Server{
		Addr:    config.DefaultPort,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}
