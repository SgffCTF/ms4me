package postgres

import (
	"context"
	"fmt"
	"ms4me/sso/internal/database"
	"ms4me/sso/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (db *Storage) User(ctx context.Context, username string) (*models.User, error) {
	const op = "database.postgres.User"

	var user models.User
	row := db.DB.QueryRow(ctx, "SELECT id, username, password FROM users WHERE username = $1", username)
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		if err == pgx.ErrNoRows {
			return &models.User{}, fmt.Errorf("%s: %w", op, database.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

// SaveUser save user to database
func (s *Storage) SaveUser(
	ctx context.Context,
	username string,
	passwordHash []byte,
) (int64, error) {
	const op = "database.postgres.SaveUser"

	var userID int64
	err := s.DB.QueryRow(ctx, "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id", username, passwordHash).Scan(&userID)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" { // unique_violation
				return 0, database.ErrUserExists
			}
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userID, nil
}
