package tests

import (
	"context"
	gamedto "ms4me/game_creator/internal/http/dto/game"
	"ms4me/game_creator/internal/utils"
	"ms4me/game_creator/tests/suite"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
)

func TestGetGames(t *testing.T) {
	const gamesCount = 5

	ctx := context.Background()
	st := suite.New()

	games := make(map[string]*gamedto.CreateGameRequest)
	var token string
	var err error
	for i := 0; i < gamesCount; i++ {
		username, password := RandomCredentials()
		token, err = st.CreateAccount(ctx, username, password)
		require.NoError(t, err)

		game := &gamedto.CreateGameRequest{
			Title: gofakeit.Word(),
			Rows:  gofakeit.Number(8, 20),
			Cols:  gofakeit.Number(8, 20),
		}
		response, err := st.CreateGame(ctx, token, game)
		require.NoError(t, err)
		games[response.ID] = game
	}
	getGames, err := st.GetGames(ctx, token, nil)
	require.NoError(t, err)
	require.Empty(t, getGames.Error)
	require.Equal(t, "OK", getGames.Status)

	for key := range games {
		consist := false
		for _, getGame := range getGames.Games {
			if getGame.ID == key {
				consist = true
				break
			}
		}
		if !consist {
			t.Errorf("get games req not contain added games")
		}
	}
}

func TestGetGame(t *testing.T) {
	ctx := context.Background()
	st := suite.New()

	username, password := RandomCredentials()
	token, err := st.CreateAccount(ctx, username, password)
	require.NoError(t, err)

	game := gamedto.CreateGameRequest{
		Title: gofakeit.Word(),
		Rows:  gofakeit.Number(8, 20),
		Cols:  gofakeit.Number(8, 20),
	}

	response, err := st.CreateGame(ctx, token, &game)
	require.NoError(t, err)

	responseGetGame, err := st.GetGame(ctx, token, response.ID)
	require.NoError(t, err)
	require.Empty(t, response.Error)
	require.Equal(t, "OK", response.Status)

	require.Len(t, responseGetGame.Games, 1)
	require.Equal(t, game.Title, responseGetGame.Games[0].Title)
	require.Equal(t, utils.MineFunc(game.Rows, game.Cols), responseGetGame.Games[0].Mines)
	require.Equal(t, game.Rows, responseGetGame.Games[0].Rows)
	require.Equal(t, game.Cols, responseGetGame.Games[0].Cols)
}
