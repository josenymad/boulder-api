package types

import "time"

type EnvVars struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name" binding:"required"`
}

type Round struct {
	ID        int       `json:"id"`
	Number    int       `json:"number" binding:"required"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}

type Competitor struct {
	ID         int    `json:"id"`
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
	CategoryID int    `json:"category_id" binding:"required"`
}

type BoulderProblem struct {
	ID      int `json:"id"`
	Number  int `json:"number" binding:"required"`
	RoundID int `json:"round_id" binding:"required"`
}

type Score struct {
	ID           int `json:"id"`
	Attempts     int `json:"attempts" binding:"required"`
	Points       int `json:"points" binding:"required"`
	CompetitorID int `json:"competitor_id" binding:"required"`
	ProblemID    int `json:"problem_id" binding:"required"`
}

type TotalScore map[string]interface{}
