package models

import (
	"time"
)

const MaxPlayers = 2

type Game struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Mines        int       `json:"mines"`
	Rows         int       `json:"rows"`
	Cols         int       `json:"cols"`
	OwnerID      int64     `json:"owner_id"`
	OwnerName    string    `json:"owner_name,omitempty"`
	IsPublic     bool      `json:"is_public"`
	CreatedAt    time.Time `json:"created_at"`
	Status       string    `json:"status"`
	PlayersCount int       `json:"players_count"`
	MaxPlayers   int       `json:"max_players"`
}

type GameDetails struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Mines        int       `json:"mines"`
	Rows         int       `json:"rows"`
	Cols         int       `json:"cols"`
	OwnerID      int64     `json:"owner_id"`
	OwnerName    string    `json:"owner_name,omitempty"`
	IsPublic     bool      `json:"is_public"`
	CreatedAt    time.Time `json:"created_at"`
	Status       string    `json:"status"`
	PlayersCount int       `json:"players_count"`
	MaxPlayers   int       `json:"max_players"`
	Players      []*User   `json:"players"`
}
