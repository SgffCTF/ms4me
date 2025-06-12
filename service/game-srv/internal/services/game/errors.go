package game

import "errors"

var (
	ErrOnlyOwnerCanStartGame = errors.New("Только создатель может начать игру")
	ErrGameAlreadyStarted    = errors.New("Игра уже начата")
	ErrGameIsNotOpen         = errors.New("Игра не открыта")
)
