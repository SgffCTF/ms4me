package storage

import "errors"

var (
	ErrEmptyRequest             = errors.New("Запрос пустой")
	ErrGameNotFoundOrNotYourOwn = errors.New("Игра не найдена или ты не являешься её владельцем")
	ErrOnlyOwnerCanStartGame    = errors.New("Только создатель может начать игру")
	ErrGameIsNotOpen            = errors.New("Игра не открыта")
	ErrGameAlreadyStarted       = errors.New("Игра уже начата")
	ErrPlayerAlreadyExists      = errors.New("Игрок уже есть")
	ErrAlreadyPlaying           = errors.New("Участвовать можно только в одной игре")
	ErrAlreadyCreatedGame       = errors.New("Создать можно только одну игру")
	ErrDeleteNotOpenGame        = errors.New("Только открытые игры могут быть удалены")
	ErrMaxPlayers               = errors.New("В игре уже участвует максимальное количество игроков")
	ErrGameNotFound             = errors.New("Игра не найдена")
	ErrOwnerCantExitFromOwnGame = errors.New("Создатель не может выйти из своей игры")
	ErrYouNotParticipate        = errors.New("Ты не участвуешь в данной игре")
	ErrIncorrectCountOfPlayers  = errors.New("Некорректное количество игроков, чтобы начать игру")
	ErrUserExists               = errors.New("пользователь уже существует")
	ErrUserNotFound             = errors.New("пользователь не найден")
)
