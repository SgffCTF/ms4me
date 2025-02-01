package redis

import (
	"bytes"
)

func GetChannel(userID string) string {
	channel := bytes.Buffer{}
	channel.WriteString("{")
	channel.WriteString(userID)
	channel.WriteString("}-channel")
	return channel.String()
}
