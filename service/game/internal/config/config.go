package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string          `yaml:"env" env:"ENV"`
	SSOConfig      *SSOConfig      `yaml:"sso" env:"-"`
	AppConfig      *HTTPConfig     `yaml:"http" env:"-"`
	DatabaseConfig *DatabaseConfig `yaml:"db" env:"-"`
}

type HTTPConfig struct {
	Host        string        `yaml:"host" env:"HOST" env-required:"true"`
	Port        string        `yaml:"port" env:"PORT" env-default:"15004"`
	Timeout     time.Duration `yaml:"timeout" env:"HTTP_TIMEOUT"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"HTTP_IDLE_TIMEOUT"`
}

type SSOConfig struct {
	Host string `yaml:"host" env:"SSO_HOST" env-required:"true"`
	Port int    `yaml:"port" env:"SSO_PORT" env-required:"true"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" env:"DB_HOST" env-default:"127.0.0.1"`
	Port     int    `yaml:"port" env:"DB_PORT" env-default:"5432"`
	Username string `yaml:"username" env:"DB_USERNAME" env-required:"true"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-required:"true"`
	Name     string `yaml:"dbname" env:"DB_NAME" env-required:"true"`
}

func MustParse() *Config {
	var configPath string
	flag.StringVar(&configPath, "config", "", "config path")
	flag.Parse()

	cfg := new(Config)
	var err error
	if configPath == "" {
		err = cleanenv.ReadEnv(cfg)
	} else {
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			panic("wrong config path")
		}
		err = cleanenv.ReadConfig(configPath, cfg)
	}
	if err != nil {
		panic("error reading config:" + err.Error())
	}

	return cfg
}
