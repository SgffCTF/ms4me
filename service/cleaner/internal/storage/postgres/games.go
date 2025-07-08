package postgres

import (
	"context"
	"time"
)

func (s *Storage) DeleteGamesBefore(ctx context.Context, t time.Time) (int64, error) {
	res, err := s.DB.Exec(ctx, "DELETE FROM games WHERE created_at < $1", t)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}
