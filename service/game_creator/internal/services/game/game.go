package game

import (
	"context"
	gamedto "ms4me/game_creator/internal/http/dto/game"
	"ms4me/game_creator/internal/models"
	"ms4me/game_creator/internal/services/batcher"
	"ms4me/game_creator/internal/utils"
	gameclient "ms4me/game_creator/pkg/http/game"

	"github.com/google/uuid"
)

type GameStorage interface {
	CreateGame(ctx context.Context, game *models.Game, userID int64) error
	GetGames(ctx context.Context, filter *gamedto.GetGamesRequest) ([]*models.Game, error)
	GetGameByID(ctx context.Context, id string, userID int64) (*models.Game, error)
	UpdateGame(ctx context.Context, id string, userID int64, game *models.Game) error
	DeleteGame(ctx context.Context, id string, userID int64) error
	StartGame(ctx context.Context, id string, userID int64) error
	EnterGame(ctx context.Context, id string, userID int64) error
	ExitGame(ctx context.Context, id string, userID int64) error
}

type Game struct {
	DB      GameStorage
	batcher *batcher.Batcher
}

func New(db GameStorage, batcher *batcher.Batcher) *Game {
	return &Game{DB: db, batcher: batcher}
}

func (g *Game) CreateGame(ctx context.Context, userID int64, game *gamedto.CreateGameRequest) (string, error) {
	id := uuid.New().String()
	newGame := &models.Game{
		ID:       id,
		Title:    game.Title,
		Mines:    utils.MineFunc(game.Rows, game.Cols),
		Rows:     game.Rows,
		Cols:     game.Cols,
		OwnerID:  userID,
		IsPublic: *game.IsPublic,
	}

	err := g.DB.CreateGame(ctx, newGame, userID)
	if err != nil {
		return "", err
	}
	if err = g.batcher.AddEvents(ctx, gameclient.Event{
		Type:   gameclient.CreateGame,
		GameID: id,
		UserID: userID,
	}); err != nil {
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
	newGame := &models.Game{
		Title:    game.Title,
		Mines:    utils.MineFunc(game.Rows, game.Cols),
		Rows:     game.Rows,
		Cols:     game.Cols,
		IsPublic: *game.IsPublic,
	}
	err := g.DB.UpdateGame(ctx, id, userID, newGame)
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

func (g *Game) StartGame(ctx context.Context, id string, userID int64) error {
	err := g.DB.StartGame(ctx, id, userID)
	if err != nil {
		return err
	}
	if err = g.batcher.AddEvents(ctx, gameclient.Event{
		Type:   gameclient.StartGame,
		GameID: id,
		UserID: userID,
	}); err != nil {
		return err
	}
	return nil
}

func (g *Game) EnterGame(ctx context.Context, id string, userID int64) error {
	err := g.DB.EnterGame(ctx, id, userID)
	if err != nil {
		return err
	}
	if err = g.batcher.AddEvents(ctx, gameclient.Event{
		Type:   gameclient.JoinGame,
		GameID: id,
		UserID: userID,
	}); err != nil {
		return err
	}
	return nil
}

func (g *Game) ExitGame(ctx context.Context, id string, userID int64) error {
	err := g.DB.ExitGame(ctx, id, userID)
	if err != nil {
		return err
	}
	return nil
}
