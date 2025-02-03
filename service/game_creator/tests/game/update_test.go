package tests

import (
	"context"
	gamedto "ms4me/game_creator/internal/http/dto/game"
	"ms4me/game_creator/internal/http/dto/response"
	"ms4me/game_creator/internal/models"
	"ms4me/game_creator/internal/storage"
	"ms4me/game_creator/tests/suite"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
)

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

	for _, tc := range testcases {
		t.Run(tc.Name, func(tt *testing.T) {
			username, password := RandomCredentials()
			token, err := st.CreateAccount(ctx, username, password)
			require.NoError(t, err)

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

func TestUpdateGame_NotYourOwn(t *testing.T) {
	ctx := context.Background()
	st := suite.New()

	username, password := RandomCredentials()
	token1, err := st.CreateAccount(ctx, username, password)
	require.NoError(t, err)
	token2, err := st.CreateAccount(ctx, username+"diff", password)
	require.NoError(t, err)

	crRes, err := st.CreateGame(ctx, token1, &gamedto.CreateGameRequest{
		Title: gofakeit.Word(),
		Rows:  gofakeit.Number(8, 20),
		Cols:  gofakeit.Number(8, 20),
	})
	require.NoError(t, err)

	upRes, err := st.UpdateGame(ctx, token2, crRes.ID, &gamedto.UpdateGameRequest{
		Title: gofakeit.Word(),
		Rows:  gofakeit.Number(8, 20),
		Cols:  gofakeit.Number(8, 20),
	})
	require.NoError(t, err)
	require.Equal(t, &response.Response{Status: "Error", Error: storage.ErrGameNotFoundOrNotYourOwn.Error()}, upRes)
}

func TestUpdateGame_NotExists(t *testing.T) {
	ctx := context.Background()
	st := suite.New()

	username, password := RandomCredentials()
	token, err := st.CreateAccount(ctx, username, password)
	require.NoError(t, err)

	upRes, err := st.UpdateGame(ctx, token, gofakeit.UUID(), &gamedto.UpdateGameRequest{
		Title: gofakeit.Word(),
		Rows:  gofakeit.Number(8, 20),
		Cols:  gofakeit.Number(8, 20),
	})
	require.NoError(t, err)
	require.Equal(t, &response.Response{Status: "Error", Error: storage.ErrGameNotFoundOrNotYourOwn.Error()}, upRes)
}
