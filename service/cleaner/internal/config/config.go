package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Env          string        `envconfig:"ENV"`
	CleanTimeout time.Duration `envconfig:"CLEAN_TIMEOUT"`
	CleanBefore  time.Duration `envconfig:"CLEAN_BEFORE"`
	*DatabaseConfig
}

type DatabaseConfig struct {
	Host     string `envconfig:"DB_HOST" env-default:"127.0.0.1"`
	Port     int    `envconfig:"DB_PORT" env-default:"5432"`
	Username string `envconfig:"DB_USERNAME" env-required:"true"`
	Password string `envconfig:"DB_PASSWORD" env-required:"true"`
	Name     string `envconfig:"DB_NAME" env-required:"true"`
}

const envconfigFilename = ".env"

func MustParseConfig() *Config {
	environment := os.Getenv("ENV")
	if environment == "" {
		environment = "local"
	}

	if environment == "local" {
		godotenv.Load(envconfigFilename)
	}

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		panic("error loading env: " + err.Error())
	}

	return &cfg
}
