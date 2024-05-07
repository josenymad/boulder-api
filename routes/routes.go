package routes

import (
	"log"
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

// POST

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

func CreateBoulderProblem(c *gin.Context) {
	var boulderProblem types.BoulderProblem
	if err := c.BindJSON(&boulderProblem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to bind boulder problem JSON", "error": err.Error()})
		return
	}

	query := "INSERT INTO boulder_problems (round_id, problem_number) VALUES ($1, $2) RETURNING problem_id"

	err := config.DB.QueryRow(query, boulderProblem.RoundID, boulderProblem.Number).Scan(&boulderProblem.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create boulder problem", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, boulderProblem)
}

func CreateScore(c *gin.Context) {
	var score types.Score
	if err := c.BindJSON(&score); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to bind score JSON", "error": err.Error()})
		return
	}

	query := "INSERT INTO scores (competitor_id, problem_id, attempts, points) VALUES ($1, $2, $3, $4) RETURNING score_id"

	err := config.DB.QueryRow(query, score.CompetitorID, score.ProblemID, score.Attempts, score.Points).Scan(&score.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create score", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, score)
}

// GET

func GetAllCategories(c *gin.Context) {
	query := "SELECT * FROM competition_categories"
	rows, err := config.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get competition categories", "error": err.Error()})
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Error closing competition category rows: %v", err)
		}
	}()

	var categories []types.Category
	for rows.Next() {
		var category types.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to scan competition category columns", "error": err.Error()})
			return
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error iterating over competition category rows", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func GetAllRounds(c *gin.Context) {
	query := "SELECT * FROM rounds"
	rows, err := config.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get rounds", "error": err.Error()})
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Error closing rounds rows: %v", err)
		}
	}()

	var rounds []types.Round
	for rows.Next() {
		var round types.Round
		if err := rows.Scan(&round.ID, &round.Number, &round.StartDate, &round.EndDate); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to scan rounds columns", "error": err.Error()})
			return
		}
		rounds = append(rounds, round)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error iterating over rounds rows", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rounds)
}

func GetAllCompetitors(c *gin.Context) {
	query := "SELECT competitor_id, name, category_id FROM competitors"
	rows, err := config.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get competitors", "error": err.Error()})
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Error closing competitor rows: %v", err)
		}
	}()

	var competitors []types.Competitor
	for rows.Next() {
		var competitor types.Competitor
		if err := rows.Scan(&competitor.ID, &competitor.Name, &competitor.CategoryID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to scan competitor columns", "error": err.Error()})
			return
		}
		competitors = append(competitors, competitor)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error iterating over competitor rows", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, competitors)
}
