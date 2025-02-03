package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Debug            bool              `yaml:"debug" env:"DEBUG" env-default:"false"`
	SSOConfig        *SSOConfig        `yaml:"sso"`
	AppConfig        *HTTPConfig       `yaml:"http"`
	CentrifugoConfig *CentrifugoConfig `yaml:"centrifugo"`
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

type CentrifugoConfig struct {
	PingInterval time.Duration `yaml:"ping_interval" env:"WS_PING_INTERVAL" env-required:"true"`
	PongTimeout  time.Duration `yaml:"pong_timeout" env:"WS_PONG_TIMEOUT" env-required:"true"`
	ExpTime      time.Duration `yaml:"expiration_time" env:"EXP_TIME" env-required:"true"`
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
