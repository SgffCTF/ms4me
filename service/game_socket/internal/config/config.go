package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Env string `envconfig:"ENV"`
	*AppConfig
	*RedisConfig
}

type AppConfig struct {
	Host          string        `envconfig:"APP_HOST"`
	Port          int           `envconfig:"APP_PORT"`
	WSPingTimeout time.Duration `envconfig:"APP_WS_PING_TIMEOUT"`
	JwtSecret     string        `envconfig:"APP_JWT_SECRET" json:"-"`
	Timeout       time.Duration `envconfig:"APP_HTTP_TIMEOUT"`
	IdleTimeout   time.Duration `envconfig:"APP_HTTP_IDLE_TIMEOUT"`
	CORSOrigins   []string      `envconfig:"APP_CORS_ORIGINS"`
	CORSMethods   []string      `envconfig:"APP_CORS_METHODS"`
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
		if err := godotenv.Load(envconfigFilename); err != nil {
			fmt.Println("could not load environment")
		}
	}

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}

	return &cfg
}
