package gamedto

type OpenCellRequest struct {
	Row int `json:"row" validate:"required,gte=0"`
	Col int `json:"col" validate:"required,gte=0"`
}
