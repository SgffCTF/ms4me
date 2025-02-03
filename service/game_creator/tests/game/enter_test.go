package tests

import (
	"context"
	"ms4me/game_creator/internal/http/dto/response"
	"ms4me/game_creator/internal/storage"
	"ms4me/game_creator/tests/suite"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnterGame(t *testing.T) {
	ctx := context.Background()
	st := suite.New()

	tokens := make([]string, 3)

	var err error
	for i := 0; i < 3; i++ {
		username, password := RandomCredentials()
		tokens[i], err = st.CreateAccount(ctx, username, password)
		require.NoError(t, err)
	}

	resp, err := st.CreateGame(ctx, tokens[0], RandomGame())
	require.NoError(t, err)

	// test enter creator
	enterRes, err := st.ControlGame(ctx, suite.ControlTypeEnter, tokens[0], resp.ID)
	require.NoError(t, err)
	require.Equal(t, response.StatusError, enterRes.Status)
	require.Equal(t, storage.ErrAlreadyPlaying.Error(), enterRes.Error)

	// test enter player1
	enterRes, err = st.ControlGame(ctx, suite.ControlTypeEnter, tokens[1], resp.ID)
	require.NoError(t, err)
	require.Equal(t, response.StatusOK, enterRes.Status)

	// test enter player1 again
	enterRes, err = st.ControlGame(ctx, suite.ControlTypeEnter, tokens[1], resp.ID)
	require.NoError(t, err)
	require.Equal(t, response.StatusError, enterRes.Status)
	require.Equal(t, storage.ErrMaxPlayers.Error(), enterRes.Error)

	// test exit player2
	exitRes, err := st.ControlGame(ctx, suite.ControlTypeExit, tokens[2], resp.ID)
	require.NoError(t, err)
	require.Equal(t, response.StatusError, exitRes.Status)
	require.Equal(t, storage.ErrYouNotParticipate.Error(), exitRes.Error)

	// test exit player1
	exitRes, err = st.ControlGame(ctx, suite.ControlTypeExit, tokens[1], resp.ID)
	require.NoError(t, err)
	require.Equal(t, response.StatusOK, exitRes.Status)
}
