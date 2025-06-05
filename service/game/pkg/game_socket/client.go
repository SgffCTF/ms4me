package gameclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"ms4me/game/internal/config"
	"ms4me/game/internal/models"
	"net/http"
	"net/url"
)

const eventsEndpoint = "/api/v1/events"

type GameClient struct {
	URL *url.URL
}

func New(cfg *config.GameSocketConfig) *GameClient {
	url, err := url.Parse(fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		panic(err)
	}
	return &GameClient{
		URL: url,
	}
}

func (c *GameClient) LoadEvents(events []models.Event) error {
	url := c.URL
	url.Path = eventsEndpoint

	client := &http.Client{}
	eventsJson, err := json.Marshal(&EventsRequest{Events: events})
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(eventsJson)

	resp, err := client.Post(c.URL.String(), "application/json", body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
