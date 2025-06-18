package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	gamedto "ms4me/game/internal/http/dto/game"
	"ms4me/game/internal/models"
	"ms4me/game/internal/storage"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Storage) CreateGame(ctx context.Context, game *models.Game, userID int64) (string, error) {
	const op = "storage.postgres.CreateGame"

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
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
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if countGames != 0 {
		return "", storage.ErrAlreadyCreatedGame
	}

	var gameID string
	err = tx.QueryRow(ctx, `
	INSERT INTO games
	(id, title, mines, rows, cols, owner_id, is_public)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id`,
		game.ID, game.Title, game.Mines, game.Rows, game.Cols,
		game.OwnerID, game.IsPublic).Scan(&gameID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	result, err := tx.Exec(ctx, "INSERT INTO players (user_id, game_id) VALUES ($1, $2)", userID, game.ID)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return "", fmt.Errorf("%s: %w", op, storage.ErrPlayerAlreadyExists)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return gameID, nil
}

func (s *Storage) GetGames(ctx context.Context, filter *gamedto.GetGamesRequest) ([]*models.Game, error) {
	const op = "storage.postgres.GetGames"

	builder := sq.Select("g.id", "title", "mines", "rows", "cols", "owner_id", "created_at", "status", "is_public", "max_players",
		"(SELECT COUNT(*) FROM players WHERE game_id = g.id) AS players_now", "u.username", "g.winner_id").
		From("games g").
		Join("users u ON u.id = g.owner_id").
		Where("is_public = true").
		OrderBy("g.created_at DESC").
		PlaceholderFormat(sq.Dollar)

	if filter.Query != "" {
		builder = builder.Where(sq.Expr("title ILIKE ?", "%"+filter.Query+"%"))
	}
	if filter.Status != "" {
		builder = builder.Where(sq.Eq{"g.status": filter.Status})
	}
	if filter.Limit > 0 {
		builder = builder.Limit(uint64(filter.Limit))
	}
	if filter.Page > 0 {
		offset := (filter.Page - 1) * filter.Limit
		builder = builder.Offset(uint64(offset))
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	rows, err := s.DB.Query(ctx, query, args...)
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
			&game.PlayersCount, &game.OwnerName, &game.WinnerID,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		games = append(games, &game)
	}

	return games, nil
}

func (s *Storage) GetGameByIDUserID(ctx context.Context, id string, userID int64) (*models.GameDetails, error) {
	const op = "storage.postgres.GetGameByIDUserID"

	row := s.DB.QueryRow(ctx, `
	SELECT 
    g.id, g.title, g.mines, g.rows, g.cols, 
    g.owner_id, g.status, g.created_at, g.is_public, g.max_players,
    COUNT(p.user_id) AS players_now,
    u.username, g.winner_id
	FROM games g
	JOIN users u ON u.id = g.owner_id
	LEFT JOIN players p ON p.game_id = g.id
	WHERE g.id = $1
	AND (g.is_public = true OR EXISTS (
		SELECT 1 FROM players WHERE game_id = g.id AND user_id = $2
	))
	GROUP BY g.id, u.username
	`, id, userID)

	var game models.GameDetails
	if err := row.Scan(
		&game.ID, &game.Title, &game.Mines, &game.Rows,
		&game.Cols, &game.OwnerID, &game.Status, &game.CreatedAt,
		&game.IsPublic, &game.MaxPlayers, &game.PlayersCount, &game.OwnerName, &game.WinnerID,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrGameNotFoundOrNotYourOwn
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	players := make([]*models.User, 0)
	rows, err := s.DB.Query(ctx, `
	SELECT u.id, u.username
	FROM users u
	JOIN players p ON p.user_id = u.id
	WHERE p.game_id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
		var player models.User
		err := rows.Scan(&player.ID, &player.Username)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		players = append(players, &player)
	}
	game.Players = players

	return &game, nil
}

func (s *Storage) GetGameByID(ctx context.Context, id string) (*models.GameDetails, error) {
	const op = "storage.postgres.GetGameByID"

	row := s.DB.QueryRow(ctx, `
	SELECT g.id, title, mines, rows, cols, owner_id, status, created_at, is_public, max_players,
	(SELECT COUNT(*) FROM players WHERE game_id = g.id) AS players_now, u.username, g.winner_id
	FROM games g
	JOIN users u ON u.id = g.owner_id
	WHERE g.id = $1`, id)

	var game models.GameDetails
	if err := row.Scan(
		&game.ID, &game.Title, &game.Mines, &game.Rows,
		&game.Cols, &game.OwnerID, &game.Status, &game.CreatedAt,
		&game.IsPublic, &game.MaxPlayers, &game.PlayersCount, &game.OwnerName, &game.WinnerID,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrGameNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	players := make([]*models.User, 0)
	rows, err := s.DB.Query(ctx, `
	SELECT u.id, u.username
	FROM users u
	JOIN players p ON p.user_id = u.id
	WHERE p.game_id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
		var player models.User
		err := rows.Scan(&player.ID, &player.Username)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		players = append(players, &player)
	}
	game.Players = players

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
	queryBuilder = queryBuilder.Set("is_public", game.IsPublic)
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

	if status == "closed" {
		return fmt.Errorf("%s: %w", op, storage.ErrDeleteClosedGame)
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

	var countGames int // Кол-во игр, в которых сейчас находится пользователь
	err = tx.QueryRow(ctx, `
	SELECT COUNT(*) FROM players p
	LEFT JOIN games g
	ON p.game_id = g.id
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

func (s *Storage) GetUserGames(ctx context.Context, userID int64) ([]*models.Game, error) {
	const op = "storage.postgres.GetUserGames"

	builder := sq.Select("g.id", "title", "mines", "rows", "cols", "owner_id", "created_at", "status", "is_public", "max_players",
		"(SELECT COUNT(*) FROM players WHERE game_id = g.id) AS players_now", "u.username", "g.winner_id").
		From("games g").
		Join("players p ON p.game_id = g.id").
		Join("users u ON u.id = g.owner_id").
		Where(sq.Eq{"p.user_id": userID}).
		OrderBy("g.created_at DESC").
		PlaceholderFormat(sq.Dollar)
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	rows, err := s.DB.Query(ctx, query, args...)
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
			&game.PlayersCount, &game.OwnerName, &game.WinnerID,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		games = append(games, &game)
	}

	return games, nil
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

func (s *Storage) UpdateGameStatus(ctx context.Context, id string, status string) error {
	const op = "storage.postgres.UpdateGameStatus"

	cmd, err := s.DB.Exec(ctx, "UPDATE games SET status = $1 WHERE id = $2", status, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return storage.ErrGameNotFound
	}

	return nil
}

func (s *Storage) UpdateWinner(ctx context.Context, id string, winnerID int64) error {
	const op = "storage.postgres.UpdateGameStatus"

	cmd, err := s.DB.Exec(ctx, "UPDATE games SET winner_id = $1 WHERE id = $2", winnerID, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return storage.ErrGameNotFound
	}

	return nil
}
