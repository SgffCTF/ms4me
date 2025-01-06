package authgrpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInternal           = status.Error(codes.Internal, "internal error")
	ErrEmptyToken         = status.Error(codes.InvalidArgument, "token is required")
	ErrExpiredToken       = status.Error(codes.Unauthenticated, "token is expired")
	ErrInvalidToken       = status.Error(codes.Unauthenticated, "token is invalid")
	ErrUserExists         = status.Error(codes.AlreadyExists, "user already exists")
	ErrInvalidCredentials = status.Error(codes.InvalidArgument, "invalid credentials")
)
