package postgres

import (
	"context"
	"ms4me/game/internal/models"
	"ms4me/game/internal/storage"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Storage) CreateUser(ctx context.Context, username string, password string) (int64, error) {
	var id int64
	err := s.DB.QueryRow(ctx, "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id", username, password).Scan(&id)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" { // unique_violation
				return 0, storage.ErrUserExists
			}
		}
		return 0, err
	}
	return id, nil
}

func (s *Storage) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := s.DB.QueryRow(ctx, "SELECT id, username, password FROM users WHERE username = $1", username).
		Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, storage.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}
