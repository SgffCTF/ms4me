package tests

import (
	"context"
	gamedto "ms4me/game_creator/internal/http/dto/game"
	"ms4me/game_creator/internal/http/dto/response"
	"ms4me/game_creator/internal/storage"
	"ms4me/game_creator/tests/suite"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
)

func TestStartGame_HappyPath(t *testing.T) {
	ctx := context.Background()
	st := suite.New()

	username, password := RandomCredentials()
	token1, err := st.CreateAccount(ctx, username, password)
	require.NoError(t, err)

	token2, err := st.CreateAccount(ctx, username+"diff", password)
	require.NoError(t, err)

	game := gamedto.CreateGameRequest{
		Title: gofakeit.Word(),
		Rows:  gofakeit.Number(8, 20),
		Cols:  gofakeit.Number(8, 20),
	}
	crRes, err := st.CreateGame(ctx, token1, &game)
	require.NoError(t, err)
	require.Equal(t, response.StatusOK, crRes.Status)

	// test start game with < 2 players
	startRes, err := st.ControlGame(ctx, suite.ControlTypeStart, token1, crRes.ID)
	require.NoError(t, err)
	require.Equal(t, response.StatusError, startRes.Status)
	require.Equal(t, storage.ErrIncorrectCountOfPlayers.Error(), startRes.Error)

	// test start game with 2 players
	enterRes, err := st.ControlGame(ctx, suite.ControlTypeEnter, token2, crRes.ID)
	require.NoError(t, err)
	require.Equal(t, response.StatusOK, enterRes.Status)

	startRes, err = st.ControlGame(ctx, suite.ControlTypeStart, token1, crRes.ID)
	require.NoError(t, err)
	require.Equal(t, response.StatusOK, startRes.Status)
}

func TestStartGame_NotYourGame(t *testing.T) {
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

	otherToken, err := st.CreateAccount(ctx, username+"diff", password)
	require.NoError(t, err)

	startRes, err := st.ControlGame(ctx, suite.ControlTypeStart, otherToken, crRes.ID)
	require.NoError(t, err)
	require.Equal(t, response.StatusError, startRes.Status)
	require.Equal(t, storage.ErrOnlyOwnerCanStartGame.Error(), startRes.Error)
}
