package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/josenymad/boulder-api/config"
	"github.com/josenymad/boulder-api/routes"
	_ "github.com/lib/pq"
)

var testFlag = flag.Bool("test", false, "Run in test mode")
var developmentFlag = flag.Bool("dev", false, "Run in development mode")

func main() {
	flag.Parse()

	if *testFlag {
		fmt.Println("Running in test mode")
		err := config.ConnectDB("test")
		if err != nil {
			log.Fatalf("Could not set up database: %v", err)
		}
	} else if *developmentFlag {
		fmt.Println("Running in development mode")
		err := config.ConnectDB("development")
		if err != nil {
			log.Fatalf("Could not set up database: %v", err)
		}
	} else {
		fmt.Println("Running in normal mode")
		err := config.ConnectDB("")
		if err != nil {
			log.Fatalf("Could not set up database: %v", err)
		}
	}

	defer config.DB.Close()

	router := gin.Default()

	router.GET("/health", routes.HealthCheckHandler)
	router.POST("/competition", routes.CreateCompetition)
	router.POST("/categories", routes.CreateCompetitionCategory)
	router.POST("/rounds", routes.CreateRound)
	router.POST("/competitors", routes.CreateCompetitor)
	router.POST("/boulder-problems", routes.CreateBoulderProblem)
	router.POST("/scores", routes.CreateScore)
	router.GET("/competitions", routes.GetAllCompetitions)
	router.GET("/categories", routes.GetAllCategories)
	router.GET("/boulder-problems/:round", routes.GetBoulderProblems)
	router.GET("/rounds/:competition", routes.GetAllRounds)
	router.GET("/competitors", routes.GetAllCompetitors)
	router.GET("/scores", routes.GetAllScores)

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
