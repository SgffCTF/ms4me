package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	gamedto "ms4me/game_creator/internal/http/dto/game"
	"ms4me/game_creator/internal/models"
	"ms4me/game_creator/internal/storage"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Storage) CreateGame(ctx context.Context, game *models.Game, userID int64) error {
	const op = "storage.postgres.CreateGame"

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				err = fmt.Errorf("rollback failed: %v, original error: %w", rollbackErr, err)
			}
		} else {
			if cErr := tx.Commit(ctx); cErr != nil {
				err = fmt.Errorf("commit failed: %v, original error: %w", cErr, err)
			}
		}
	}()

	var countGames int
	err = tx.QueryRow(ctx, `
		SELECT COUNT(*) 
		FROM players p
		LEFT JOIN games g
		ON g.owner_id = p.user_id
		WHERE g.status != 'closed' AND p.user_id = $1`, userID,
	).Scan(&countGames)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if countGames != 0 {
		return storage.ErrAlreadyPlaying
	}

	result, err := tx.Exec(ctx, "INSERT INTO games (id, title, mines, rows, cols, owner_id, is_public) VALUES ($1, $2, $3, $4, $5, $6, $7)", game.ID, game.Title, game.Mines, game.Rows, game.Cols, game.OwnerID, game.IsPublic)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("%s: no rows affected", op)
	}

	result, err = tx.Exec(ctx, "INSERT INTO players (user_id, game_id) VALUES ($1, $2)", userID, game.ID)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, storage.ErrPlayerAlreadyExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetGames(ctx context.Context, filter *gamedto.GetGamesRequest) ([]*models.Game, error) {
	const op = "storage.postgres.GetGames"

	var rows pgx.Rows
	var err error
	if filter.Limit == 0 && filter.Page == 0 {
		rows, err = s.DB.Query(ctx, `
			SELECT id, title, mines, rows, cols, owner_id, created_at, status, is_public, max_players,
			(SELECT COUNT(*) FROM players WHERE game_id = g.id) AS players_now
			FROM games g
			WHERE is_public = true`,
		)
	} else {
		offset := (filter.Page - 1) * filter.Limit
		rows, err = s.DB.Query(ctx, `
			SELECT id, title, mines, rows, cols, owner_id, created_at, status, is_public, max_players,
			(SELECT COUNT(*) FROM players WHERE game_id = g.id) AS players_now
			FROM games g
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
		if err := rows.Scan(
			&game.ID, &game.Title, &game.Mines, &game.Rows,
			&game.Cols, &game.OwnerID, &game.CreatedAt,
			&game.Status, &game.IsPublic, &game.MaxPlayers,
			&game.Players,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		games = append(games, &game)
	}

	return games, nil
}

func (s *Storage) GetGameByID(ctx context.Context, id string, userID int64) (*models.Game, error) {
	const op = "storage.postgres.GetGameByID"

	row := s.DB.QueryRow(ctx, `
	SELECT id, title, mines, rows, cols, owner_id, status, created_at, is_public, max_players,
	(SELECT COUNT(*) FROM players WHERE game_id = g.id) AS players_now
	FROM games g
	WHERE id = $1 AND (is_public = true OR owner_id = $2)`, id, userID)

	var game models.Game
	if err := row.Scan(
		&game.ID, &game.Title, &game.Mines, &game.Rows,
		&game.Cols, &game.OwnerID, &game.Status, &game.CreatedAt,
		&game.IsPublic, &game.MaxPlayers, &game.Players,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrGameNotFoundOrNotYourOwn
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &game, nil
}

func (s *Storage) UpdateGame(ctx context.Context, id string, userID int64, game *models.Game) error {
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
	queryBuilder.Set("is_public", game.IsPublic)
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

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				err = fmt.Errorf("rollback failed: %v, original error: %w", rollbackErr, err)
			}
		} else {
			if cErr := tx.Commit(ctx); cErr != nil {
				err = fmt.Errorf("commit failed: %v, original error: %w", cErr, err)
			}
		}
	}()

	var status string
	err = tx.QueryRow(ctx, "SELECT status FROM games WHERE id = $1 AND owner_id = $2", id, userID).Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", op, storage.ErrGameNotFoundOrNotYourOwn)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if status != "open" {
		return fmt.Errorf("%s: %w", op, storage.ErrDeleteNotOpenGame)
	}

	result, err := tx.Exec(ctx, "DELETE FROM players WHERE game_id = $1", id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrGameNotFoundOrNotYourOwn)
	}

	result, err = tx.Exec(ctx, "DELETE FROM games WHERE id = $1 AND owner_id = $2", id, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrGameNotFoundOrNotYourOwn)
	}

	return nil
}

func (s *Storage) StartGame(ctx context.Context, id string, userID int64) error {
	const op = "storage.postgres.StartGame"

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				err = fmt.Errorf("rollback failed: %v, original error: %w", rollbackErr, err)
			}
		} else {
			if cErr := tx.Commit(ctx); cErr != nil {
				err = fmt.Errorf("commit failed: %v, original error: %w", cErr, err)
			}
		}
	}()

	var ownerID int64
	var status string
	err = tx.QueryRow(ctx, "SELECT owner_id, status FROM games WHERE id = $1", id).Scan(&ownerID, &status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", op, storage.ErrGameNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if ownerID != userID {
		return fmt.Errorf("%s: %w", op, storage.ErrOnlyOwnerCanStartGame)
	}
	if status == "started" {
		return fmt.Errorf("%s: %w", op, storage.ErrGameAlreadyStarted)
	}
	if status != "open" {
		return fmt.Errorf("%s: %w", op, storage.ErrGameIsNotOpen)
	}

	var countPlayers int
	err = tx.QueryRow(ctx, "SELECT COUNT(*) FROM players WHERE game_id = $1", id).Scan(&countPlayers)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if countPlayers != models.MaxPlayers {
		return fmt.Errorf("%s: %w", op, storage.ErrIncorrectCountOfPlayers)
	}

	result, err := tx.Exec(ctx, "UPDATE games SET status = 'started' WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) EnterGame(ctx context.Context, id string, userID int64) error {
	const op = "storage.postgres.EnterGame"

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				err = fmt.Errorf("rollback failed: %v, original error: %w", rollbackErr, err)
			}
		} else {
			if cErr := tx.Commit(ctx); cErr != nil {
				err = fmt.Errorf("commit failed: %v, original error: %w", cErr, err)
			}
		}
	}()

	var status string
	err = tx.QueryRow(ctx, "SELECT status FROM games WHERE id = $1", id).Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", op, storage.ErrGameNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	if status != "open" {
		return fmt.Errorf("%s: %w", op, storage.ErrGameIsNotOpen)
	}

	var countGames int // count games where player gaming
	err = tx.QueryRow(ctx, `
	SELECT COUNT(*) FROM players p
	LEFT JOIN games g
	ON p.user_id = g.owner_id
	WHERE g.status != 'closed' AND p.user_id = $1`, userID,
	).Scan(&countGames)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if countGames != 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrAlreadyPlaying)
	}

	var countPlayers int
	err = tx.QueryRow(ctx, "SELECT COUNT(*) FROM players WHERE game_id = $1", id).Scan(&countPlayers)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if countPlayers >= models.MaxPlayers {
		return fmt.Errorf("%s: %w", op, storage.ErrMaxPlayers)
	}

	result, err := tx.Exec(ctx, "INSERT INTO players (game_id, user_id) VALUES ($1, $2)", id, userID)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, storage.ErrPlayerAlreadyExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) ExitGame(ctx context.Context, id string, userID int64) error {
	const op = "storage.postgres.ExitGame"

	var ownerID int64
	err := s.DB.QueryRow(ctx, "SELECT owner_id FROM games WHERE id = $1", id).Scan(&ownerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", op, storage.ErrGameNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if ownerID == userID {
		return fmt.Errorf("%s: %w", op, storage.ErrOwnerCantExitFromOwnGame)
	}

	result, err := s.DB.Exec(ctx, "DELETE FROM players WHERE game_id = $1 AND user_id = $2", id, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrYouNotParticipate)
	}

	return nil
}
