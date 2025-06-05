package dto_ws

import "time"

type MessageRequest struct {
	ChatID    string    `json:"chat_id"`
	CreatorID int64     `json:"-"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
