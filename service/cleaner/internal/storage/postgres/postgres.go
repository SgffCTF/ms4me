package postgres

import (
	"context"
	"fmt"
	"ms4me/cleaner/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	DB *pgxpool.Pool
}

func (s *Storage) Stop() {
	s.DB.Close()
}

func New(ctx context.Context, config *config.DatabaseConfig) *Storage {
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", config.Username, config.Password, config.Host, config.Port, config.Name)

	db, err := pgxpool.New(ctx, url)
	if err != nil {
		panic("Failed to create connection pool: " + err.Error())
	}

	err = db.Ping(ctx)
	if err != nil {
		panic("Failed to ping database: " + err.Error())
	}

	return &Storage{db}
}
