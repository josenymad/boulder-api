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
	ID   string `json:"id"`
	Name string `json:"name" binding:"required"`
}

type Round struct {
	ID        string    `json:"id"`
	Number    string    `json:"number" binding:"required"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}

type Competitor struct {
	ID         string `json:"id"`
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
	CategoryID string `json:"category_id" binding:"required"`
}
