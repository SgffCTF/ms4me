package game

import (
	"context"
	gamedto "game-creator/internal/http/dto/game"
	"game-creator/internal/models"

	"github.com/google/uuid"
)

type GameStorage interface {
	CreateGame(ctx context.Context, game *models.Game) error
	GetGames(ctx context.Context, filter *gamedto.GetGamesRequest) ([]*models.Game, error)
	GetGameByID(ctx context.Context, id string, userID int64) (*models.Game, error)
	UpdateGame(ctx context.Context, id string, userID int64, game *gamedto.UpdateGameRequest) error
	DeleteGame(ctx context.Context, id string, userID int64) error
}

type Game struct {
	DB GameStorage
}

func New(db GameStorage) *Game {
	return &Game{DB: db}
}

func (g *Game) CreateGame(ctx context.Context, userID int64, game *gamedto.CreateGameRequest) (string, error) {
	id := uuid.New().String()
	newGame := &models.Game{
		ID:       id,
		Title:    game.Title,
		Mines:    game.Mines,
		Rows:     game.Rows,
		Cols:     game.Cols,
		OwnerID:  userID,
		IsPublic: *game.IsPublic,
	}

	err := g.DB.CreateGame(ctx, newGame)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (g *Game) GetGames(ctx context.Context, filter *gamedto.GetGamesRequest) ([]*models.Game, error) {
	games, err := g.DB.GetGames(ctx, filter)
	if err != nil {
		return nil, err
	}
	return games, nil
}

func (g *Game) GetGame(ctx context.Context, id string, userID int64) (*models.Game, error) {
	game, err := g.DB.GetGameByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (g *Game) UpdateGame(ctx context.Context, id string, userID int64, game *gamedto.UpdateGameRequest) error {
	err := g.DB.UpdateGame(ctx, id, userID, game)
	if err != nil {
		return err
	}
	return nil
}

func (g *Game) DeleteGame(ctx context.Context, id string, userID int64) error {
	err := g.DB.DeleteGame(ctx, id, userID)
	if err != nil {
		return err
	}
	return nil
}
