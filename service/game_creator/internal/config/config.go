package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string             `yaml:"env" env:"ENV"`
	AppConfig  *ApplicationConfig `yaml:"app" env:"-"`
	DBConfig   *DatabaseConfig    `yaml:"db" env:"-"`
	SSOConfig  *SSOConfig         `yaml:"sso" env:"-"`
	GameConfig *GameConfig        `yaml:"game" env:"-"`
}

type ApplicationConfig struct {
	Host        string        `yaml:"host" env:"HOST"`
	Port        int           `yaml:"port" env:"PORT"`
	Timeout     time.Duration `yaml:"timeout" env:"TIMEOUT"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"IDLE_TIMEOUT"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" env:"DB_HOST" env-default:"127.0.0.1"`
	Port     int    `yaml:"port" env:"DB_PORT" env-default:"5432"`
	Username string `yaml:"username" env:"DB_USERNAME" env-required:"true"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-required:"true"`
	Name     string `yaml:"dbname" env:"DB_NAME" env-required:"true"`
}

type SSOConfig struct {
	Host string `yaml:"host" env:"SSO_HOST" env-required:"true"`
	Port int    `yaml:"port" env:"SSO_PORT" env-required:"true"`
}

type GameConfig struct {
	Host string `yaml:"host" env:"GAME_HOST" env-required:"true"`
	Port int    `yaml:"port" env:"GAME_PORT" env-required:"true"`
}

func MustParseConfig() *Config {
	var configPath string
	flag.StringVar(&configPath, "config", "", "config path")
	flag.Parse()

	cfg := &Config{
		AppConfig:  &ApplicationConfig{},
		DBConfig:   &DatabaseConfig{},
		SSOConfig:  &SSOConfig{},
		GameConfig: &GameConfig{},
	}
	var err error
	if configPath == "" {
		err = cleanenv.ReadEnv(cfg.GameConfig)
		err = cleanenv.ReadEnv(cfg.SSOConfig)
		err = cleanenv.ReadEnv(cfg.DBConfig)
		err = cleanenv.ReadEnv(cfg.AppConfig)
		err = cleanenv.ReadEnv(cfg)
	} else {
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			panic("wrong config path")
		}
		cfg, err = Parse(configPath)
	}
	if err != nil {
		panic("error reading config:" + err.Error())
	}

	return cfg
}

func Parse(path string) (*Config, error) {
	var config Config
	if err := cleanenv.ReadConfig(path, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
