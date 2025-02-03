package events

import "errors"

var (
	ErrChannelNotFound      = errors.New("channel not found")
	ErrUserAlreadyInChannel = errors.New("user already in channel")
)
