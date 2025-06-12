package handlers

import (
	"context"
	"log/slog"
	"ms4me/game/internal/config"
	gamedto "ms4me/game/internal/http/dto/game"
	"ms4me/game/internal/models"
)

type GameService interface {
	CreateGame(ctx context.Context, userID int64, game *gamedto.CreateGameRequest) (string, error)
	GetGames(ctx context.Context, filter *gamedto.GetGamesRequest) ([]*models.Game, error)
	GetGame(ctx context.Context, id string, userID int64) (*models.GameDetails, error)
	UpdateGame(ctx context.Context, id string, userID int64, game *gamedto.UpdateGameRequest) error
	DeleteGame(ctx context.Context, id string, userID int64) error
	StartGame(ctx context.Context, id string, userID int64) error
	EnterGame(ctx context.Context, id string, userID int64, username string) error
	ExitGame(ctx context.Context, id string, userID int64, username string) error
	UserGames(ctx context.Context, userID int64) ([]*models.Game, error)
	OpenCell(ctx context.Context, req *gamedto.OpenCellRequest, gameID string, userID int64) error
	GameStarted(ctx context.Context, gameID string) (bool, error)
	CloseGame(ctx context.Context, gameID string) error
}

type AuthService interface {
	Register(ctx context.Context, username, password string) (int64, error)
	Login(ctx context.Context, username, password string) (string, error)
}

type GameHandlers struct {
	log     *slog.Logger
	gameSrv GameService
	authSrv AuthService
	cfg     *config.Config
}

func New(log *slog.Logger, gameSrv GameService, authSrv AuthService, cfg *config.Config) *GameHandlers {
	return &GameHandlers{
		log:     log,
		gameSrv: gameSrv,
		authSrv: authSrv,
		cfg:     cfg,
	}
}
