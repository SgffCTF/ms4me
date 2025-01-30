package tests

import (
	"context"
	gamedto "game-creator/internal/http/dto/game"
	"game-creator/tests/suite"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
)

func TestCreateGame(t *testing.T) {
	ctx := context.Background()
	st := suite.New()
	defer st.App.Stop()

	username, password := RandomCredentials()
	token, err := st.CreateAccount(ctx, username, password)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	isPublic := true
	response, err := st.CreateGame(ctx, token, &gamedto.CreateGameRequest{
		Title:    gofakeit.Word(),
		Mines:    gofakeit.Number(5, 20),
		Rows:     gofakeit.Number(10, 20),
		Cols:     gofakeit.Number(10, 20),
		IsPublic: &isPublic,
	})
	require.NoError(t, err)
	require.Empty(t, response.Error)
	require.Equal(t, "OK", response.Status)
}

func TestGetGames(t *testing.T) {
	const gamesCount = 5

	ctx := context.Background()
	st := suite.New()
	defer st.App.Stop()

	username, password := RandomCredentials()
	token, err := st.CreateAccount(ctx, username, password)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	games := make(map[string]*gamedto.CreateGameRequest)
	for i := 0; i < gamesCount; i++ {
		game := &gamedto.CreateGameRequest{
			Title: gofakeit.Word(),
			Mines: gofakeit.Number(5, 20),
			Rows:  gofakeit.Number(10, 20),
			Cols:  gofakeit.Number(10, 20),
		}
		response, err := st.CreateGame(ctx, token, game)
		require.NoError(t, err)
		require.Empty(t, response.Error)
		require.Equal(t, "OK", response.Status)
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
	defer st.App.Stop()

	username, password := RandomCredentials()
	token, err := st.CreateAccount(ctx, username, password)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	response, err := st.CreateGame(ctx, token, &gamedto.CreateGameRequest{
		Title: gofakeit.Word(),
		Mines: gofakeit.Number(5, 20),
		Rows:  gofakeit.Number(10, 20),
		Cols:  gofakeit.Number(10, 20),
	})
	require.NoError(t, err)
	require.Empty(t, response.Error)
	require.Equal(t, "OK", response.Status)

}
