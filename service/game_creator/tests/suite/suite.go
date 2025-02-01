package suite

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"game-creator/internal/config"
	gamedto "game-creator/internal/http/dto/game"
	"game-creator/internal/http/dto/response"
	grpcclient "game-creator/pkg/grpc/client"
	ssov1 "game-creator/pkg/grpc/sso"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type Suite struct {
	Client    *http.Client
	SSOClient *grpcclient.SSOClient
	URL       *url.URL
}

func New() *Suite {
	cfg, err := config.Parse("test_config.yml")
	if err != nil {
		panic(err)
	}

	ssoClient := grpcclient.New(cfg.SSOConfig)

	httpClient := http.Client{}
	url, err := url.Parse(fmt.Sprintf("http://%s:%d", cfg.AppConfig.Host, cfg.AppConfig.Port))
	if err != nil {
		panic(err)
	}

	suite := &Suite{Client: &httpClient, SSOClient: ssoClient, URL: url}

	return suite
}

func (s *Suite) CreateAccount(ctx context.Context, username, password string) (string, error) {
	_, err := s.SSOClient.AuthClient.Register(ctx, &ssov1.RegisterRequest{Username: username, Password: password})
	if err != nil {
		return "", err
	}
	loginResponse, err := s.SSOClient.AuthClient.Login(ctx, &ssov1.LoginRequest{Username: username, Password: password})
	if err != nil {
		return "", err
	}
	return loginResponse.GetToken(), nil
}

func (s *Suite) CreateGame(ctx context.Context, token string, dto *gamedto.CreateGameRequest) (*gamedto.CreateGameResponse, error) {
	url := s.URL
	url.Path = "/api/v1/game"
	req, err := http.NewRequestWithContext(ctx, "POST", url.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")
	jsonBody, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	req.Body = io.NopCloser(bytes.NewBuffer(jsonBody))
	req.ContentLength = int64(len(jsonBody))

	res, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	response := new(gamedto.CreateGameResponse)
	if err := json.Unmarshal(body, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Suite) GetGames(ctx context.Context, token string, dto *gamedto.GetGamesRequest) (*gamedto.GetGamesResponse, error) {
	url := s.URL
	url.Path = "/api/v1/game"

	if dto != nil { // set get params
		query := url.Query()
		query.Set("limit", strconv.Itoa(dto.Limit))
		query.Set("page", strconv.Itoa(dto.Page))
		url.RawQuery = query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", token)

	res, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	response := new(gamedto.GetGamesResponse)
	if err := json.Unmarshal(body, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Suite) GetGame(ctx context.Context, token string, id string) (*gamedto.GetGamesResponse, error) {
	url := s.URL
	url.Path = "/api/v1/game/" + id
	req, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", token)

	res, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	response := new(gamedto.GetGamesResponse)
	if err := json.Unmarshal(body, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Suite) DeleteGame(ctx context.Context, token, id string) (*response.Response, error) {
	url := s.URL
	url.Path = "/api/v1/game/" + id
	req, err := http.NewRequestWithContext(ctx, "DELETE", url.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", token)

	res, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	response := new(response.Response)
	if err := json.Unmarshal(body, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Suite) UpdateGame(ctx context.Context, token, id string, dto *gamedto.UpdateGameRequest) (*response.Response, error) {
	url := s.URL
	url.Path = "/api/v1/game/" + id
	req, err := http.NewRequestWithContext(ctx, "PUT", url.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")
	jsonBody, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	req.Body = io.NopCloser(bytes.NewBuffer(jsonBody))
	req.ContentLength = int64(len(jsonBody))

	res, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	response := new(response.Response)
	if err := json.Unmarshal(body, response); err != nil {
		return nil, err
	}

	return response, nil
}
