package tests

import (
	"context"
	gamedto "ms4me/game_creator/internal/http/dto/game"
	"ms4me/game_creator/internal/storage"
	"ms4me/game_creator/tests/suite"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
)

func TestCreateGame_HappyPath(t *testing.T) {
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

func TestCreateGame_Double(t *testing.T) {
	ctx := context.Background()
	st := suite.New()

	username, password := RandomCredentials()
	token, err := st.CreateAccount(ctx, username, password)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	response, err := st.CreateGame(ctx, token, &gamedto.CreateGameRequest{
		Title: gofakeit.Word(),
		Rows:  gofakeit.Number(8, 20),
		Cols:  gofakeit.Number(8, 20),
	})
	require.NoError(t, err)
	require.Empty(t, response.Error)
	require.Equal(t, "OK", response.Status)

	response, err = st.CreateGame(ctx, token, &gamedto.CreateGameRequest{
		Title: gofakeit.Word(),
		Rows:  gofakeit.Number(8, 20),
		Cols:  gofakeit.Number(8, 20),
	})
	require.NoError(t, err)
	require.Equal(t, storage.ErrAlreadyPlaying.Error(), response.Error)
	require.Equal(t, "Error", response.Status)
}
