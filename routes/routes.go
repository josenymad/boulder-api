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

	query := `INSERT INTO competition_categories (name) VALUE ($1) RETURNING id`

	err := config.DB.QueryRow(query, category.Name).Scan(&category.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create competition category", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}
