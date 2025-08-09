package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Bot
	Postgres
	Service
	Jira
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

type Jira struct {
	Username string `env:"JARVIS_BOT_EMAIL"`
	Password string `env:"JARVIS_BOT_PASSWORD"`
	BaseUrl  string `env:"JIRA_BASE_URL"`
}

type Service struct {
	Port string `env:"JARVIS_BOT_SERVICE_PORT"`
}

func MustLoadConfig() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Fatalf("error while read environments: %s", err.Error())
	}
	return cfg
}
