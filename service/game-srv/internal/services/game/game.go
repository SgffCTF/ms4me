package game

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	gamedto "ms4me/game/internal/http/dto/game"
	"ms4me/game/internal/models"
	"ms4me/game/internal/storage/redis"
	ingameclient "ms4me/game/pkg/ingame_client"

	"github.com/google/uuid"
	"github.com/jacute/prettylogger"
)

const defaultGameRows = 8
const defaultGameCols = 8
const defaultGameMines = 10

const GAME_STARTED_STATUS = "started"
const GAME_CLOSED_STATUS = "closed"
const GAME_OPEN_STATUS = "open"

type GameStorage interface {
	CreateGame(ctx context.Context, game *models.Game, userID int64) (string, error)
	GetGames(ctx context.Context, filter *gamedto.GetGamesRequest) ([]*models.Game, error)
	GetGameByID(ctx context.Context, id string) (*models.GameDetails, error)
	GetGameByIDUserID(ctx context.Context, id string, userID int64) (*models.GameDetails, error)
	UpdateGame(ctx context.Context, id string, userID int64, game *models.Game) error
	DeleteGame(ctx context.Context, id string, userID int64) error
	StartGame(ctx context.Context, id string, userID int64) error
	EnterGame(ctx context.Context, id string, userID int64) error
	ExitGame(ctx context.Context, id string, userID int64) error
	GetUserGames(ctx context.Context, userID int64) ([]*models.Game, error)
	UpdateGameStatus(ctx context.Context, id string, status string) error
	UpdateWinner(ctx context.Context, id string, winnerID int64) error
}

type Game struct {
	log *slog.Logger
	DB  GameStorage
	rdb *redis.Redis
	gc  *ingameclient.IngameClient
}

func New(log *slog.Logger, db GameStorage, rdb *redis.Redis, gc *ingameclient.IngameClient) *Game {
	return &Game{log: log, DB: db, rdb: rdb, gc: gc}
}

func (g *Game) CreateGame(ctx context.Context, userID int64, game *gamedto.CreateGameRequest) (string, error) {
	const op = "game.CreateGame"
	log := g.log.With(slog.String("op", op), slog.Int64("user_id", userID))
	id := uuid.New().String()
	newGame := &models.Game{
		ID:       id,
		Title:    game.Title,
		Mines:    defaultGameMines,
		Rows:     defaultGameRows,
		Cols:     defaultGameCols,
		OwnerID:  userID,
		IsPublic: *game.IsPublic,
	}

	_, err := g.DB.CreateGame(ctx, newGame, userID)
	if err != nil {
		log.Error("error creating game", prettylogger.Err(err))
		return "", err
	}
	createdGame, err := g.DB.GetGameByID(ctx, id)
	if err != nil {
		log.Error("error got game", prettylogger.Err(err))
		return "", err
	}
	gameMarshalled, err := json.Marshal(createdGame)
	if err != nil {
		log.Error("error marshalling game", prettylogger.Err(err))
		return "", err
	}
	if err = g.rdb.PublishEvent(ctx, models.Event{
		Type:     models.TypeCreateGame,
		GameID:   id,
		UserID:   userID,
		Username: createdGame.OwnerName,
		IsPublic: createdGame.IsPublic,
		Payload:  gameMarshalled,
	}); err != nil {
		log.Error("error adding create game event", prettylogger.Err(err))
		return "", err
	}

	log.Info("game created successfully", slog.String("game_id", id))

	return id, nil
}

func (g *Game) GetGames(ctx context.Context, filter *gamedto.GetGamesRequest) ([]*models.Game, error) {
	const op = "game.GetGames"
	log := g.log.With(slog.String("op", op))
	games, err := g.DB.GetGames(ctx, filter)
	if err != nil {
		log.Error("error getting games", prettylogger.Err(err))
		return nil, err
	}

	log.Info("games got successfully")

	return games, nil
}

func (g *Game) GetGame(ctx context.Context, id string) (*models.GameDetails, error) {
	const op = "game.GetGame"
	log := g.log.With(slog.String("op", op), slog.String("game_id", id))
	game, err := g.DB.GetGameByID(ctx, id)
	if err != nil {
		log.Error("error getting game", prettylogger.Err(err))
		return nil, err
	}

	log.Info("game got successfully")

	return game, nil
}

func (g *Game) UpdateGame(ctx context.Context, id string, userID int64, game *gamedto.UpdateGameRequest) error {
	const op = "game.UpdateGame"
	log := g.log.With(slog.String("op", op), slog.String("id", id), slog.Int64("user_id", userID))
	newGame := &models.Game{
		Title:    game.Title,
		Mines:    defaultGameMines,
		Rows:     defaultGameRows,
		Cols:     defaultGameCols,
		IsPublic: *game.IsPublic,
	}
	gameBeforeUpdate, err := g.DB.GetGameByID(ctx, id)
	if err != nil {
		log.Error("error got game", prettylogger.Err(err))
	}
	err = g.DB.UpdateGame(ctx, id, userID, newGame)
	if err != nil {
		log.Error("error updating game", prettylogger.Err(err))
		return err
	}
	gameMarshalled, err := json.Marshal(game)
	if err != nil {
		log.Error("error marshalling game", prettylogger.Err(err))
		return err
	}
	if err = g.rdb.PublishEvent(ctx, models.Event{
		Type:     models.TypeUpdateGame,
		GameID:   id,
		UserID:   userID,
		IsPublic: gameBeforeUpdate.IsPublic, // Отправляем isPublic, который был ещё до апдейта
		Payload:  gameMarshalled,
	}); err != nil {
		log.Error("error pushing event", slog.String("event_type", "update_game"), prettylogger.Err(err))
		return err
	}
	log.Info("game updated successfully")
	return nil
}

func (g *Game) DeleteGame(ctx context.Context, id string, userID int64) error {
	const op = "game.DeleteGame"
	log := g.log.With(slog.String("op", op), slog.String("game_id", id), slog.Int64("user_id", userID))
	game, err := g.DB.GetGameByID(ctx, id)
	if err != nil {
		log.Error("error got game", prettylogger.Err(err))
	}
	err = g.DB.DeleteGame(ctx, id, userID)
	if err != nil {
		log.Error("error deleting game", prettylogger.Err(err))
		return err
	}
	if err = g.rdb.PublishEvent(ctx, models.Event{
		Type:     models.TypeDeleteGame,
		GameID:   id,
		IsPublic: game.IsPublic,
		UserID:   userID,
	}); err != nil {
		log.Error("error pushing event", slog.String("event_type", "delete_game"), prettylogger.Err(err))
		return err
	}
	log.Info("game deleted successfully")
	return nil
}

func (g *Game) StartGame(ctx context.Context, id string, userID int64) error {
	const op = "game.StartGame"
	log := g.log.With(slog.String("op", op), slog.String("game_id", id), slog.Int64("user_id", userID))

	err := g.gc.Ready(id)
	if err != nil {
		if errors.Is(err, ingameclient.ErrNotReady) {
			log.Warn("participants not ready to start game")
			return err
		}
		log.Error("error starting game", prettylogger.Err(err))
		return err
	}
	game, err := g.DB.GetGameByIDUserID(ctx, id, userID)
	if err != nil {
		log.Error("error getting game", prettylogger.Err(err))
		return err
	}
	if game.OwnerID != userID {
		log.Info("only owner can start the game")
		return fmt.Errorf("%s: %w", op, ErrOnlyOwnerCanStartGame)
	}
	if game.Status == GAME_STARTED_STATUS {
		log.Info("game already started")
		return fmt.Errorf("%s: %w", op, ErrGameAlreadyStarted)
	}
	if game.Status != GAME_OPEN_STATUS {
		log.Info("game is not open")
		return fmt.Errorf("%s: %w", op, ErrGameIsNotOpen)
	}
	err = g.DB.StartGame(ctx, id, userID)
	if err != nil {
		log.Error("error starting game", prettylogger.Err(err))
		return err
	}
	if err = g.rdb.PublishEvent(ctx, models.Event{
		Type:     models.TypeStartGame,
		GameID:   id,
		IsPublic: game.IsPublic,
		UserID:   userID,
	}); err != nil {
		log.Error("error pushing event", slog.String("event_type", "start_game"), prettylogger.Err(err))
		return err
	}
	log.Info("game started successfully")
	return nil
}

func (g *Game) EnterGame(ctx context.Context, id string, userID int64, username string) error {
	const op = "game.EnterGame"
	log := g.log.With(slog.String("op", op), slog.String("game_id", id), slog.Int64("user_id", userID))
	game, err := g.DB.GetGameByID(ctx, id)
	if err != nil {
		log.Error("error getting game", prettylogger.Err(err))
		return err
	}
	err = g.DB.EnterGame(ctx, id, userID)
	if err != nil {
		log.Error("error entering game", prettylogger.Err(err))
		return err
	}
	if game.Status != GAME_OPEN_STATUS {
		log.Info("game is not open")
		return fmt.Errorf("%s: %w", op, ErrGameIsNotOpen)
	}
	if err = g.rdb.PublishEvent(ctx, models.Event{
		Type:     models.TypeJoinGame,
		GameID:   id,
		UserID:   userID,
		IsPublic: game.IsPublic,
		Username: username,
	}); err != nil {
		log.Error("error pushing event", slog.String("event_type", "enter_game"), prettylogger.Err(err))
		return err
	}
	log.Info("enter in game successfully")
	return nil
}

func (g *Game) ExitGame(ctx context.Context, id string, userID int64, username string) error {
	const op = "game.ExitGame"
	log := g.log.With(slog.String("op", op), slog.String("game_id", id), slog.Int64("user_id", userID))
	game, err := g.DB.GetGameByID(ctx, id)
	if err != nil {
		log.Error("error getting game", prettylogger.Err(err))
		return err
	}
	err = g.DB.ExitGame(ctx, id, userID)
	if err != nil {
		log.Error("error exiting game", prettylogger.Err(err))
		return err
	}
	if err = g.rdb.PublishEvent(ctx, models.Event{
		Type:     models.TypeExitGame,
		GameID:   id,
		UserID:   userID,
		Username: username,
		IsPublic: game.IsPublic,
	}); err != nil {
		log.Error("error pushing event", slog.String("event_type", "exit_game"), prettylogger.Err(err))
		return err
	}
	log.Info("exit from game successfully")
	return nil
}

func (g *Game) UserGames(ctx context.Context, userID int64) ([]*models.Game, error) {
	const op = "game.ExitGame"
	log := g.log.With(slog.String("op", op), slog.Int64("user_id", userID))
	games, err := g.DB.GetUserGames(ctx, userID)
	if err != nil {
		log.Error("error got user games", prettylogger.Err(err))
		return nil, err
	}
	log.Info("user games got successfully")
	return games, nil
}

func (g *Game) GetGameStatus(ctx context.Context, gameID string) (string, error) {
	const op = "game.GameStarted"
	log := g.log.With(slog.String("op", op), slog.String("game_id", gameID))

	game, err := g.DB.GetGameByID(ctx, gameID)
	if err != nil {
		log.Error("error getting game", prettylogger.Err(err))
		return "", err
	}
	log.Info("game status got successfully", slog.String("status", game.Status))
	return game.Status, nil
}

func (g *Game) CloseGame(ctx context.Context, gameID string, winnerID int64) error {
	const op = "game.CloseGame"
	log := g.log.With(slog.String("op", op), slog.String("game_id", gameID))

	err := g.DB.UpdateGameStatus(ctx, gameID, GAME_CLOSED_STATUS)
	if err != nil {
		log.Error("error closing game", prettylogger.Err(err))
		return err
	}
	err = g.DB.UpdateWinner(ctx, gameID, winnerID)
	if err != nil {
		log.Error("error updating winner of game", prettylogger.Err(err))
		return err
	}
	log.Info("game closed successfully")
	return nil
}
