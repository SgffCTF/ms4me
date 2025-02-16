package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (s *Storage) GetGameIDByUserID(ctx context.Context, userID int64) (string, error) {
	const op = "storage.postgres.GetGameIDByUserID"

	var gameID string
	err := s.DB.QueryRow(ctx, "SELECT id FROM games WHERE owner_id = $1", userID).Scan(&gameID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, ErrGameNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return gameID, nil
}
