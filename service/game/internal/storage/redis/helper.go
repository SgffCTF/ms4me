package redis

import (
	"bytes"
)

func GetChannel(gameID string) string {
	channel := bytes.Buffer{}
	channel.WriteString("{")
	channel.WriteString(gameID)
	channel.WriteString("}-channel")
	return channel.String()
}
