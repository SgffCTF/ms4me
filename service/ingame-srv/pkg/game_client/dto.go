package gameclient

import "ms4me/game_socket/internal/http/dto"

type GameStartedResponse struct {
	dto.Response
	Started bool `json:"started"`
}
