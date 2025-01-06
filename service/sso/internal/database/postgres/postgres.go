package postgres

import (
	"context"
	"fmt"
	"ms4me/sso/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	DB *pgxpool.Pool
}

func (s *Storage) Stop() {
	s.DB.Close()
}

func New(config *config.DatabaseConfig) *Storage {
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", config.Username, config.Password, config.Host, config.Port, config.Name)

	db, err := pgxpool.New(context.Background(), url)
	if err != nil {
		panic("Failed to create connection pool: " + err.Error())
	}

	err = db.Ping(context.Background())
	if err != nil {
		panic("Failed to ping database: " + err.Error())
	}

	return &Storage{db}
}
