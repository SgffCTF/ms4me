package channel

import "sync"

type Channels struct {
	Channels map[int64]string
	Mu       sync.Mutex
}

func New() *Channels {
	return &Channels{
		Channels: make(map[int64]string),
		Mu:       sync.Mutex{},
	}
}
