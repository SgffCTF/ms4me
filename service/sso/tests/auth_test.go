package tests

import (
	authgrpc "ms4me/sso/internal/grpc/auth"
	ssov1 "ms4me/sso/internal/grpc/proto/sso"
	"ms4me/sso/internal/lib/jwt"
	"ms4me/sso/internal/models"
	"ms4me/sso/tests/suite"
	"sync"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

const (
	defaultPasswordLen = 16
	defaultUsernameLen = 16
)

type AuthTestCase struct {
	name          string
	username      string
	password      string
	expectedError string
}

type VerifyTokenTestCase struct {
	name          string
	token         string
	expectedError error
}

func TestAuth_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)
	defer st.App.GRPCApp.Stop()

	username := gofakeit.LetterN(defaultUsernameLen)
	password := gofakeit.LetterN(defaultPasswordLen)

	regRes, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Username: username,
		Password: password,
	})
	require.NoError(t, err)
	require.NotEmpty(t, regRes.GetId())

	loginRes, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Username: username,
		Password: password,
	})
	require.NoError(t, err)
	token := loginRes.GetToken()
	require.NotEmpty(t, token)

	verifyToken, err := st.AuthClient.VerifyToken(ctx, &ssov1.VerifyTokenRequest{
		Token: token,
	})
	user := verifyToken.GetUser()
	require.NoError(t, err)
	require.NotEmpty(t, user.Id)
	require.NotEmpty(t, user.Username)
}

func TestAuth_WrongCredentials(t *testing.T) {
	ctx, st := suite.New(t)
	defer st.App.GRPCApp.Stop()

	username := gofakeit.LetterN(defaultUsernameLen)
	password := gofakeit.LetterN(defaultPasswordLen)

	regRes, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Username: username,
		Password: password,
	})
	require.NoError(t, err)
	require.NotEmpty(t, regRes.GetId())

	loginRes, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Username: username,
		Password: gofakeit.LetterN(16),
	})
	require.ErrorContains(t, err, "invalid credentials")
	token := loginRes.GetToken()
	require.Empty(t, token)
}

func TestRegister_Errors(t *testing.T) {
	ctx, st := suite.New(t)
	defer st.App.GRPCApp.Stop()

	cases := []AuthTestCase{
		{
			name:          "short username",
			username:      gofakeit.LetterN(2),
			password:      gofakeit.LetterN(defaultPasswordLen),
			expectedError: "username must be at least 3 characters long",
		},
		{
			name:          "short password",
			username:      gofakeit.LetterN(defaultUsernameLen),
			password:      gofakeit.LetterN(4),
			expectedError: "password must be at least 8 characters long",
		},
		{
			name:          "empty username",
			username:      "",
			password:      gofakeit.LetterN(defaultPasswordLen),
			expectedError: "username cannot be empty",
		},
		{
			name:          "empty password",
			username:      gofakeit.LetterN(defaultUsernameLen),
			password:      "",
			expectedError: "password cannot be empty",
		},
	}

	var wg sync.WaitGroup
	for _, c := range cases {
		wg.Add(1)
		go func() {
			t.Run(c.name, func(tt *testing.T) {
				_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
					Username: c.username,
					Password: c.password,
				})
				if c.expectedError == "" {
					require.NoError(tt, err)
				} else {
					require.ErrorContains(tt, err, c.expectedError)
				}
			})
			wg.Done()
		}()
	}

	wg.Wait()
}

func TestVerifyToken(t *testing.T) {
	ctx, st := suite.New(t)
	defer st.App.GRPCApp.Stop()

	expiredToken, err := jwt.NewToken(&models.User{
		ID:       1,
		Username: gofakeit.LetterN(defaultUsernameLen),
	}, []byte(st.Cfg.AppConfig.JwtSecret), 0)
	if err != nil {
		panic("error creating jwt token: " + err.Error())
	}

	cases := []VerifyTokenTestCase{
		{
			name:          "empty token",
			token:         "",
			expectedError: authgrpc.ErrEmptyToken,
		},
		{
			name:          "invalid token",
			token:         gofakeit.LetterN(16),
			expectedError: authgrpc.ErrInvalidToken,
		},
		{
			name:          "expired token",
			token:         expiredToken,
			expectedError: authgrpc.ErrExpiredToken,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			_, err := st.AuthClient.VerifyToken(ctx,
				&ssov1.VerifyTokenRequest{
					Token: c.token,
				})
			require.ErrorIs(tt, err, c.expectedError)
		})
	}
}
