package tests

import (
	"context"
	gamedto "ms4me/game_creator/internal/http/dto/game"
	"ms4me/game_creator/tests/suite"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
)

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
