package game

import (
	"context"
	"encoding/json"
	"log/slog"
	gamedto "ms4me/game/internal/http/dto/game"
	"ms4me/game/internal/models"

	"github.com/jacute/prettylogger"
)

func (g *Game) OpenCell(ctx context.Context, req *gamedto.OpenCellRequest, gameID string, userID int64) error {
	const op = "game.OpenCell"
	log := g.log.With(slog.String("op", op))

	game, err := g.DB.GetGameByIDUserID(ctx, gameID, userID)
	if err != nil {
		log.Error("error getting game", prettylogger.Err(err))
		return err
	}
	reqMarshalled, err := json.Marshal(req)
	if err != nil {
		log.Error("error marshalling req", prettylogger.Err(err))
		return err
	}
	if err = g.batcher.AddEvents(ctx, models.Event{
		Type:     models.TypeOpenCell,
		UserID:   userID,
		GameID:   gameID,
		IsPublic: game.IsPublic,
		Payload:  reqMarshalled,
	}); err != nil {
		log.Error("error pushing event", slog.String("event_type", "open_cell"), prettylogger.Err(err))
		return err
	}
	return nil
}
