package dto

import (
	"encoding/json"
)

type ClickCellRequest struct {
	Row int `json:"row" validate:"gte=0,lte=8"`
	Col int `json:"col" validate:"gte=0,lte=8"`
}

type GetParticipantsResponse struct {
	Response
	Participants json.RawMessage `json:"participants"`
}
