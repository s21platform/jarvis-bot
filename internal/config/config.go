package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Bot
}

type Bot struct {
	Port  string `env:"JARVIS_PORT"`
	Token string `env:"JARVIS_TOKEN"`
}

func MustLoadConfig() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Fatalf("error while read environments: %s", err.Error())
	}
	return cfg
}
