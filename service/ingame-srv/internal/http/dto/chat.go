package dto

import "ms4me/game_socket/internal/models"

type CreateMessageRequest struct {
	Text string `json:"text" validate:"required,max=256"`
}

type ReadMessagesResponse struct {
	Response
	Messages []*models.Message `json:"messages"`
}
