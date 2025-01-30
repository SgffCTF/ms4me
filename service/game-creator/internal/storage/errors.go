package storage

import "errors"

var (
	ErrEmptyRequest             = errors.New("request is empty")
	ErrGameNotFoundOrNotYourOwn = errors.New("game not found or you aren't owner")
)
