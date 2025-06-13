package models

import "time"

type Message struct {
	ID              string    `json:"id"`
	CreatorID       int64     `json:"creator_id"`
	CreatorUsername string    `json:"creator_username"`
	Text            string    `json:"text"`
	CreatedAt       time.Time `json:"created_at"`
}
