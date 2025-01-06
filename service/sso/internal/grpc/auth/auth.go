package authgrpc

import (
	"context"
	"errors"
	"ms4me/sso/internal/lib/jwt"
	"ms4me/sso/internal/lib/validators"
	"ms4me/sso/internal/models"
	"ms4me/sso/internal/services/auth"

	ssov1 "ms4me/sso/internal/grpc/proto/sso"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(
		ctx context.Context,
		username string,
		password string,
	) (token string, err error)
	Register(
		ctx context.Context,
		username string,
		password string,
	) (userID int64, err error)
	VerifyToken(
		ctx context.Context,
		token string,
	) (user *models.User, err error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	username := req.GetUsername()
	password := req.GetPassword()

	err := validators.ValidateUsername(username)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = validators.ValidatePassword(password)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := s.auth.Login(ctx, username, password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, ErrInvalidCredentials
		}
		return nil, ErrInternal
	}

	return &ssov1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	username := req.GetUsername()
	password := req.GetPassword()

	err := validators.ValidateUsername(username)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = validators.ValidatePassword(password)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userID, err := s.auth.Register(ctx, username, password)
	if err != nil {
		if errors.Is(err, auth.ErrUserExists) {
			return nil, ErrUserExists
		}
		return nil, ErrInternal
	}

	return &ssov1.RegisterResponse{Id: userID}, nil
}

func (s *serverAPI) VerifyToken(ctx context.Context, req *ssov1.VerifyTokenRequest) (*ssov1.VerifyTokenResponse, error) {
	token := req.GetToken()

	if token == "" {
		return nil, ErrEmptyToken
	}

	user, err := s.auth.VerifyToken(ctx, token)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	return &ssov1.VerifyTokenResponse{User: &ssov1.User{Id: user.ID, Username: user.Username}}, nil
}
