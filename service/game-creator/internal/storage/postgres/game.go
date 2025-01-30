package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	gamedto "game-creator/internal/http/dto/game"
	"game-creator/internal/models"
	"game-creator/internal/storage"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) CreateGame(ctx context.Context, game *models.Game) error {
	const op = "storage.postgres.CreateGame"

	result, err := s.DB.Exec(ctx, "INSERT INTO games (id, title, mines, rows, cols, owner_id, is_public) VALUES ($1, $2, $3, $4, $5, $6, $7)", game.ID, game.Title, game.Mines, game.Rows, game.Cols, game.OwnerID, game.IsPublic)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("%s: no rows affected", op)
	}

	return nil
}

func (s *Storage) GetGames(ctx context.Context, filter *gamedto.GetGamesRequest) ([]*models.Game, error) {
	const op = "storage.postgres.GetGames"

	var rows pgx.Rows
	var err error
	if filter.Limit == 0 && filter.Page == 0 {
		rows, err = s.DB.Query(ctx, `
			SELECT id, title, mines, rows, cols, owner_id, is_open, created_at
			FROM games
			WHERE is_public = true`,
		)
	} else {
		offset := (filter.Page - 1) * filter.Limit
		rows, err = s.DB.Query(ctx, `
			SELECT id, title, mines, rows, cols, owner_id, is_open, created_at
			FROM games
			WHERE is_public = true
			LIMIT $1 OFFSET $2`, filter.Limit, offset,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	games := make([]*models.Game, 0)
	for rows.Next() {
		var game models.Game
		if err := rows.Scan(&game.ID, &game.Title, &game.Mines, &game.Rows, &game.Cols, &game.OwnerID, &game.IsOpen, &game.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		games = append(games, &game)
	}

	return games, nil
}

func (s *Storage) GetGameByID(ctx context.Context, id string, userID int64) (*models.Game, error) {
	const op = "storage.postgres.GetGameByID"

	row := s.DB.QueryRow(ctx, "SELECT id, title, mines, rows, cols, owner_id, is_open, created_at FROM games WHERE id = $1 AND (is_public = true OR owner_id = $2)", id, userID)

	var game models.Game
	if err := row.Scan(&game.ID, &game.Title, &game.Mines, &game.Rows, &game.Cols, &game.OwnerID, &game.IsOpen, &game.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrGameNotFoundOrNotYourOwn
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &game, nil
}

func (s *Storage) UpdateGame(ctx context.Context, id string, userID int64, game *gamedto.UpdateGameRequest) error {
	const op = "storage.postgres.UpdateGame"

	queryBuilder := sq.Update("games").Where(sq.Eq{"id": id, "owner_id": userID}).PlaceholderFormat(sq.Dollar)
	if game.Title != "" {
		queryBuilder = queryBuilder.Set("title", game.Title)
	}
	if game.Cols != 0 {
		queryBuilder = queryBuilder.Set("cols", game.Cols)
	}
	if game.Rows != 0 {
		queryBuilder = queryBuilder.Set("rows", game.Rows)
	}
	if game.Mines != 0 {
		queryBuilder = queryBuilder.Set("mines", game.Mines)
	}
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	result, err := s.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrGameNotFoundOrNotYourOwn)
	}

	return nil
}

func (s *Storage) DeleteGame(ctx context.Context, id string, userID int64) error {
	const op = "storage.postgres.DeleteGame"

	result, err := s.DB.Exec(ctx, "DELETE FROM games WHERE id = $1 AND owner_id = $2", id, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrGameNotFoundOrNotYourOwn)
	}

	return nil
}
