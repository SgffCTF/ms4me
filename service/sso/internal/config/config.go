package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	AppConfig *ApplicationConfig `yaml:"app"`
	DBConfig  *DatabaseConfig    `yaml:"db"`
}

type ApplicationConfig struct {
	Host      string        `yaml:"host"`
	Port      int           `yaml:"port"`
	JwtSecret string        `yaml:"jwt_secret" env-required:"true"`
	JwtTTL    time.Duration `yaml:"jwt_TTL" env-default:"48h"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" env-default:"127.0.0.1:"`
	Port     int    `yaml:"port" env-default:"5432"`
	Username string `yaml:"username" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	Name     string `yaml:"name" env-required:"true"`
}

func MustParseConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.Parse()

	if configPath == "" {
		panic("config path not provided")
	}

	config, err := ParseConfig(configPath)
	if err != nil {
		panic(err)
	}

	return config
}

func ParseConfig(path string) (*Config, error) {
	var config Config
	if err := cleanenv.ReadConfig(path, &config); err != nil {
		return nil, fmt.Errorf("error reading config file: %s", err.Error())
	}

	return &config, nil
}
