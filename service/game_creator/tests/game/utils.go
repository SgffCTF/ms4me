package tests

import (
	gamedto "ms4me/game_creator/internal/http/dto/game"

	"github.com/brianvoe/gofakeit"
)

func RandomCredentials() (string, string) {
	var username string
	for len(username) < 8 {
		username = gofakeit.Username()
	}
	password := gofakeit.Password(true, true, true, false, false, 8)

	return username, password
}

func RandomGame() *gamedto.CreateGameRequest {
	return &gamedto.CreateGameRequest{
		Title: gofakeit.Word(),
		Rows:  gofakeit.Number(8, 20),
		Cols:  gofakeit.Number(8, 20),
	}
}
