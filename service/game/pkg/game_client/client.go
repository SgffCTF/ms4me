package gameclient

import (
	"errors"
	"fmt"
	"ms4me/game/internal/config"
	"net/http"
	"net/url"
)

const gameSocketReadyEndpoint = "/api/v1/game/%s/ready"

var (
	ErrNotReady = errors.New("Не все участники игры готовы")
)

type GameSocketClient struct {
	URL *url.URL
}

func New(cfg *config.GameSocketConfig) *GameSocketClient {
	url, err := url.Parse(fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		panic(err)
	}
	return &GameSocketClient{
		URL: url,
	}
}

func (c *GameSocketClient) Ready(gameID string) error {
	url := c.URL
	url.Path = fmt.Sprintf(gameSocketReadyEndpoint, gameID)

	client := &http.Client{}

	resp, err := client.Get(c.URL.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooEarly {
		return ErrNotReady
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
