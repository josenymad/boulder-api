package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/josenymad/boulder-api/config"
	"github.com/josenymad/boulder-api/types"
)

func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Service is healthy",
	})
}

func CreateCompetitionCategory(c *gin.Context) {
	var category types.Category
	if err := c.BindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to bind category JSON", "error": err.Error()})
		return
	}

	query := "INSERT INTO competition_categories (name) VALUES ($1) RETURNING category_id"

	err := config.DB.QueryRow(query, category.Name).Scan(&category.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create competition category", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

func CreateRound(c *gin.Context) {
	var round types.Round
	if err := c.BindJSON(&round); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to bind round JSON", "error": err.Error()})
		return
	}

	query := "INSERT INTO rounds (round_number, start_date, end_date) VALUES ($1, $2, $3) RETURNING round_id"

	err := config.DB.QueryRow(query, round.Number, round.StartDate, round.EndDate).Scan(&round.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create round", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, round)
}

func CreateCompetitor(c *gin.Context) {
	var competitor types.Competitor
	if err := c.BindJSON(&competitor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to bind competitor JSON", "error": err.Error()})
		return
	}

	query := "INSERT INTO competitors (name, email, password, category_id) VALUES ($1, $2, $3, $4) RETURNING competitor_id"

	err := config.DB.QueryRow(query, competitor.Name, competitor.Email, competitor.Password, competitor.CategoryID).Scan(&competitor.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create competitor", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, competitor)
}
