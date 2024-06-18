package routes

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/josenymad/boulder-api/config"
	"github.com/josenymad/boulder-api/types"
	"github.com/josenymad/boulder-api/utils"
	"golang.org/x/crypto/bcrypt"
)

func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Service is healthy",
	})
}

// POST

func CreateCompetition(c *gin.Context) {
	var competition types.Competition
	if err := c.BindJSON(&competition); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to bind competition JSON", "error": err.Error()})
		return
	}

	query := `INSERT INTO competitions (competition_name) VALUES ($1) RETURNING competition_id`

	err := config.DB.QueryRow(query, competition.Name).Scan(&competition.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create competition", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, competition)
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

	query := "INSERT INTO rounds (round_number, start_date, end_date, competition_id) VALUES ($1, $2, $3, $4) RETURNING round_id"

	err := config.DB.QueryRow(query, round.Number, round.StartDate, round.EndDate, round.CompetitionID).Scan(&round.ID)
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(competitor.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password", "error": err.Error()})
		return
	}

	competitor.Password = string(hashedPassword)

	query := "INSERT INTO competitors (name, email, password, category_id) VALUES ($1, $2, $3, $4) RETURNING competitor_id"

	err = config.DB.QueryRow(query, competitor.Name, competitor.Email, competitor.Password, competitor.CategoryID).Scan(&competitor.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create competitor", "error": err.Error()})
		return
	}

	response := types.CompetitorResponse{
		ID:         competitor.ID,
		Name:       competitor.Name,
		CategoryID: competitor.CategoryID,
	}

	c.JSON(http.StatusCreated, response)
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

func GetAllCompetitions(c *gin.Context) {
	query := "SELECT * FROM competitions"
	rows, err := config.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get competitions", "error": err.Error()})
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Error closing competition rows: %v", err)
		}
	}()

	var competitions []types.Competition
	for rows.Next() {
		var competition types.Competition
		if err := rows.Scan(&competition.ID, &competition.Name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to scan competition rows", "error": err.Error()})
			return
		}
		competitions = append(competitions, competition)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error iterating over competition rows", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, competitions)
}

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
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to scan competition category rows", "error": err.Error()})
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

func GetBoulderProblems(c *gin.Context) {
	round := c.Param("round")
	query := "SELECT problem_id, problem_number FROM boulder_problems WHERE round_id = $1"

	rows, err := config.DB.Query(query, round)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get boulder problems", "error": err.Error()})
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Error closing boulder problem rows: %v", err)
		}
	}()

	var boulderProblems []types.BoulderProblem
	for rows.Next() {
		var boulderProblem types.BoulderProblem
		if err := rows.Scan(&boulderProblem.ID, &boulderProblem.Number); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to scan boulder problem rows", "error": err.Error()})
			return
		}
		boulderProblems = append(boulderProblems, boulderProblem)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error iterating over boulder problem rows", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, boulderProblems)
}

func GetAllRounds(c *gin.Context) {
	competition := c.Param("competition")
	query := "SELECT round_id, round_number, start_date, end_date FROM rounds WHERE competition_id = $1"
	rows, err := config.DB.Query(query, competition)
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
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to scan rounds rows", "error": err.Error()})
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
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to scan competitor rows", "error": err.Error()})
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

func GetAllScores(c *gin.Context) {
	category := c.Query("category")
	competition := c.Query("competition")
	query, err := utils.BuildScoresQueryString(category, competition)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error building scores query string", "error": err.Error()})
		return
	}

	rows, err := config.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get scores", "error": err.Error()})
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Error closing scores rows: %v", err)
		}
	}()

	var totalScores []types.TotalScore
	columns, err := rows.Columns()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error getting total scores columns", "error": err.Error()})
		return
	}
	values := make([]interface{}, len(columns))
	for rows.Next() {
		for i := range values {
			values[i] = new(interface{})
		}

		err := rows.Scan(values...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to scan total scores row", "error": err.Error()})
			return
		}

		totalScore := make(types.TotalScore)
		for i, column := range columns {
			val := *(values[i].(*interface{}))
			totalScore[column] = val
		}
		totalScores = append(totalScores, totalScore)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error iterating over total scores rows", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, totalScores)
}
