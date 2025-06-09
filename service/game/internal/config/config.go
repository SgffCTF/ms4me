package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Env string `envconfig:"ENV"`
	*ApplicationConfig
	*DatabaseConfig
	*GameSocketConfig
	*RedisConfig
}

type ApplicationConfig struct {
	Host        string        `envconfig:"HOST"`
	Port        int           `envconfig:"PORT"`
	Timeout     time.Duration `envconfig:"TIMEOUT"`
	IdleTimeout time.Duration `envconfig:"IDLE_TIMEOUT"`
	JwtSecret   string        `envconfig:"JWT_SECRET"`
	JwtTTL      time.Duration `envconfig:"JWT_TTL"`
	CORSOrigins []string      `envconfig:"CORS_ORIGINS"`
	CORSMethods []string      `envconfig:"CORS_METHODS"`
}

type DatabaseConfig struct {
	Host     string `envconfig:"DB_HOST" env-default:"127.0.0.1"`
	Port     int    `envconfig:"DB_PORT" env-default:"5432"`
	Username string `envconfig:"DB_USERNAME" env-required:"true"`
	Password string `envconfig:"DB_PASSWORD" env-required:"true"`
	Name     string `envconfig:"DB_NAME" env-required:"true"`
}

type GameSocketConfig struct {
	Host string `envconfig:"GAME_SOCKET_HOST" env-required:"true"`
	Port int    `envconfig:"GAME_SOCKET_PORT" env-required:"true"`
}

type RedisConfig struct {
	Host     string `envconfig:"REDIS_HOST"`
	Port     int    `envconfig:"REDIS_PORT"`
	Password string `envconfig:"REDIS_PASSWORD"`
	DB       int    `envconfig:"REDIS_DB"`
	Username string `envconfig:"REDIS_USERNAME"`
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
