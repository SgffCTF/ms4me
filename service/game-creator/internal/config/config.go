package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	AppConfig *ApplicationConfig `yaml:"app"`
	DBConfig  *DatabaseConfig    `yaml:"db"`
	SSOConfig *SSOConfig         `yaml:"sso"`
}

type ApplicationConfig struct {
	Host        string        `yaml:"host"`
	Port        int           `yaml:"port"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" env-default:"127.0.0.1"`
	Port     int    `yaml:"port" env-default:"5432"`
	Username string `yaml:"username" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	Name     string `yaml:"dbname" env-required:"true"`
}

type SSOConfig struct {
	Host string `yaml:"host" env-required:"true"`
	Port int    `yaml:"port" env-required:"true"`
}

func MustParseConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	flag.StringVar(&configPath, "config", "", "path to config file")

	if configPath == "" {
		flag.Parse()
		if configPath == "" {
			panic("config path not provided")
		}
	}

	config, err := Parse(configPath)
	if err != nil {
		panic("error parsing config: " + err.Error())
	}

	return config
}

func Parse(path string) (*Config, error) {
	var config Config
	if err := cleanenv.ReadConfig(path, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
