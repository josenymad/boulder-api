package routes_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/josenymad/boulder-api/config"
	"github.com/josenymad/boulder-api/routes"
	"github.com/josenymad/boulder-api/types"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setUpRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/competition", routes.CreateCompetition)
	router.POST("/categories", routes.CreateCompetitionCategory)
	router.POST("/rounds", routes.CreateRound)
	router.POST("/competitors", routes.CreateCompetitor)
	return router
}

func TestCreateCompetition(t *testing.T) {
	router := setUpRouter()
	err := config.ConnectDB("test")
	if err != nil {
		log.Fatalf("Could not set up database: %v", err)
	}

	competition := types.Competition{
		Name: "Test Competition",
	}
	body, err := json.Marshal(competition)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest("POST", "/competition", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	assert.Equal(t, "Test Competition", response["name"])
}

func TestCreateCompetitionCategory(t *testing.T) {
	router := setUpRouter()
	err := config.ConnectDB("test")
	if err != nil {
		log.Fatalf("Could not set up database: %v", err)
	}

	category := types.Category{
		Name: "Test Category",
	}
	body, err := json.Marshal(category)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest("POST", "/categories", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	assert.Equal(t, "Test Category", response["name"])
}

func TestCreateRound(t *testing.T) {
	router := setUpRouter()
	err := config.ConnectDB("test")
	if err != nil {
		log.Fatalf("Could not set up database: %v", err)
	}

	round := types.Round{
		Number:        1,
		StartDate:     time.Date(2024, time.July, 1, 10, 0, 0, 0, time.UTC),
		EndDate:       time.Date(2024, time.July, 15, 19, 0, 0, 0, time.UTC),
		CompetitionID: 15,
	}
	body, err := json.Marshal(round)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest("POST", "/rounds", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	assert.Equal(t, float64(1), response["number"])
	assert.Equal(t, "2024-07-01T10:00:00Z", response["start_date"])
	assert.Equal(t, "2024-07-15T19:00:00Z", response["end_date"])
	assert.Equal(t, float64(15), response["competition_id"])
}

func TestCreateCompetitor(t *testing.T) {
	router := setUpRouter()
	err := config.ConnectDB("test")
	if err != nil {
		log.Fatalf("Could not set up database: %v", err)
	}

	competitor := types.Competitor{
		Name:       "Test Competitor",
		Email:      "test@gmail.com",
		Password:   "test_password",
		CategoryID: 7,
	}
	body, err := json.Marshal(competitor)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest("POST", "/competitors", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	assert.Equal(t, "Test Competitor", response["name"])
	assert.Equal(t, float64(7), response["category_id"])
}
