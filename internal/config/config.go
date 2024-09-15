package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Bot
	Postgres
}

type Bot struct {
	Port  string `env:"JARVIS_PORT"`
	Token string `env:"JARVIS_TOKEN"`
	Url   string `env:"JARVIS_URL"`
}

type Postgres struct {
	User     string `env:"JARVIS_BOT_POSTGRES_USER"`
	Password string `env:"JARVIS_BOT_POSTGRES_PASSWORD"`
	Database string `env:"JARVIS_BOT_POSTGRES_DB"`
	Host     string `env:"JARVIS_BOT_POSTGRES_HOST"`
	Port     string `env:"JARVIS_BOT_POSTGRES_PORT"`
}

func MustLoadConfig() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Fatalf("error while read environments: %s", err.Error())
	}
	return cfg
}
