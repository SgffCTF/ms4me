package tests

import (
	"context"
	gamedto "game-creator/internal/http/dto/game"
	"game-creator/internal/http/dto/response"
	"game-creator/internal/models"
	"game-creator/tests/suite"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
)

func TestCreateGame(t *testing.T) {
	ctx := context.Background()
	st := suite.New()

	username, password := RandomCredentials()
	token, err := st.CreateAccount(ctx, username, password)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	isPublic := true
	response, err := st.CreateGame(ctx, token, &gamedto.CreateGameRequest{
		Title:    gofakeit.Word(),
		Rows:     gofakeit.Number(8, 20),
		Cols:     gofakeit.Number(8, 20),
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

	username, password := RandomCredentials()
	token, err := st.CreateAccount(ctx, username, password)
	require.NoError(t, err)

	games := make(map[string]*gamedto.CreateGameRequest)
	for i := 0; i < gamesCount; i++ {
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
	require.NotEmpty(t, token)

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
	require.Equal(t, game.Rows*game.Cols/10, responseGetGame.Games[0].Mines)
	require.Equal(t, game.Rows, responseGetGame.Games[0].Rows)
	require.Equal(t, game.Cols, responseGetGame.Games[0].Cols)
}

func TestDeleteGame(t *testing.T) {
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

	crRes, err := st.CreateGame(ctx, token, &game)
	require.NoError(t, err)

	delRes, err := st.DeleteGame(ctx, token, crRes.ID)
	require.NoError(t, err)
	require.Empty(t, delRes.Error)
	require.Equal(t, "OK", delRes.Status)

	getRes, err := st.GetGame(ctx, token, crRes.ID)
	require.NoError(t, err)
	require.Equal(t, "game not found or you aren't owner", getRes.Error)
	require.Equal(t, "Error", getRes.Status)
}

func TestUpdateGame_HappyPath(t *testing.T) {
	testcases := []struct {
		Name            string
		GameToUpdate    *gamedto.CreateGameRequest
		Request         *gamedto.UpdateGameRequest
		GameAfterUpdate *models.Game
		Response        *response.Response
	}{
		{
			Name: "update title",
			GameToUpdate: &gamedto.CreateGameRequest{
				Title: "123",
				Rows:  11,
				Cols:  12,
			},
			Request: &gamedto.UpdateGameRequest{
				Title: gofakeit.Word(),
			},
			Response: &response.Response{
				Status: "OK",
			},
		},
		{
			Name: "update all",
			GameToUpdate: &gamedto.CreateGameRequest{
				Title: ".",
				Rows:  8,
				Cols:  8,
			},
			Request: &gamedto.UpdateGameRequest{
				Title: gofakeit.Word(),
				Rows:  gofakeit.Number(8, 20),
				Cols:  gofakeit.Number(8, 20),
			},
			Response: &response.Response{
				Status: "OK",
			},
		},
	}
	for i, tc := range testcases {
		tc.GameAfterUpdate = &models.Game{
			Title: tc.GameToUpdate.Title,
			Rows:  tc.GameToUpdate.Rows,
			Cols:  tc.GameToUpdate.Cols,
		}
		if tc.Request.Rows != 0 {
			tc.GameAfterUpdate.Rows = tc.Request.Rows
		}
		if tc.Request.Cols != 0 {
			tc.GameAfterUpdate.Cols = tc.Request.Cols
		}
		if tc.Request.Title != "" {
			tc.GameAfterUpdate.Title = tc.Request.Title
		}
		testcases[i] = tc
	}

	ctx := context.Background()
	st := suite.New()

	username, password := RandomCredentials()
	token, err := st.CreateAccount(ctx, username, password)
	require.NoError(t, err)

	for _, tc := range testcases {
		t.Run(tc.Name, func(tt *testing.T) {
			crRes, err := st.CreateGame(ctx, token, tc.GameToUpdate)
			require.NoError(tt, err)

			upRes, err := st.UpdateGame(ctx, token, crRes.ID, tc.Request)
			require.NoError(tt, err)
			require.Equal(tt, tc.Response, upRes)

			getRes, err := st.GetGame(ctx, token, crRes.ID)
			require.NoError(tt, err)
			require.Len(tt, getRes.Games, 1)
			require.Equal(tt, tc.GameAfterUpdate.Cols, getRes.Games[0].Cols)
			require.Equal(tt, tc.GameAfterUpdate.Rows, getRes.Games[0].Rows)
			require.Equal(tt, tc.GameAfterUpdate.Title, getRes.Games[0].Title)
		})
	}
}
