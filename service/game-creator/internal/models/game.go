package models

import "time"

type Game struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Mines     int       `json:"mines"`
	Rows      int       `json:"rows"`
	Cols      int       `json:"cols"`
	OwnerID   int64     `json:"owner_id"`
	IsOpen    bool      `json:"is_open"`
	IsPublic  bool      `json:"is_public"`
	CreatedAt time.Time `json:"created_at"`
}
