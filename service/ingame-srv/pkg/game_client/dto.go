package gameclient

import "ms4me/game_socket/internal/http/dto"

type GameStatusResponse struct {
	dto.Response
	Result string `json:"result"`
}
