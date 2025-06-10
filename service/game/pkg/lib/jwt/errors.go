package jwt

import "errors"

var (
	ErrTokenExpired            = errors.New("token expired")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrTokenInvalid            = errors.New("invalid token")
)
