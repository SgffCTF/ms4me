package storage

import "errors"

var (
	ErrEmptyRequest             = errors.New("request is empty")
	ErrGameNotFoundOrNotYourOwn = errors.New("game not found or you aren't owner")
	ErrOnlyOwnerCanStartGame    = errors.New("only owner can start game")
	ErrGameIsNotOpen            = errors.New("game is not open")
	ErrGameAlreadyStarted       = errors.New("game already started")
	ErrPlayerAlreadyExists      = errors.New("player already in this game")
	ErrAlreadyPlaying           = errors.New("you can create and participate in one game in one time")
	ErrDeleteNotOpenGame        = errors.New("only open games can be deleted")
	ErrMaxPlayers               = errors.New("game has maximum players")
	ErrGameNotFound             = errors.New("game not found")
	ErrOwnerCantExitFromOwnGame = errors.New("owner can't exit from own game")
	ErrYouNotParticipate        = errors.New("you not participate in this game")
	ErrIncorrectCountOfPlayers  = errors.New("incorrect count of players to start the game")
	ErrUserExists               = errors.New("пользователь уже существует")
	ErrUserNotFound             = errors.New("пользователь не найден")
)
