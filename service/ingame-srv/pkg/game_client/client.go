package gameclient

import (
	"bytes"
	"errors"
	"fmt"
	"ms4me/game_socket/internal/config"
	"ms4me/game_socket/internal/http/dto"
	"net/http"
	"net/url"

	"github.com/go-chi/render"
)

const gameStatusEndpoint = "/api/v1/internal/game/%s/status"
const gameCloseEndpoint = "/api/v1/internal/game/%s/close"

type GameClient struct {
	URL *url.URL
}

func New(cfg *config.GameConfig) *GameClient {
	url, err := url.Parse(fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		panic(err)
	}
	return &GameClient{
		URL: url,
	}
}

func (c *GameClient) GetStatus(gameID string) (string, error) {
	url := c.URL
	url.Path = fmt.Sprintf(gameStatusEndpoint, gameID)

	client := &http.Client{}

	resp, err := client.Get(c.URL.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res GameStatusResponse
	if err := render.DecodeJSON(resp.Body, &res); err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if res.Status == dto.StatusError {
		return "", errors.New(res.Error)
	}

	return res.Result, nil
}

func (c *GameClient) Close(gameID string, winnerID int64) error {
	url := c.URL
	url.Path = fmt.Sprintf(gameCloseEndpoint, gameID)

	client := &http.Client{}

	resp, err := client.Post(c.URL.String(), "application/json", bytes.NewBuffer([]byte(fmt.Sprintf("{\"winner_id\":%d}", winnerID))))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var res dto.Response
	if err := render.DecodeJSON(resp.Body, &res); err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if res.Status == dto.StatusError {
		return errors.New(res.Error)
	}

	return nil
}
