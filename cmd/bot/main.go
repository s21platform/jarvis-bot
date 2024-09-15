package main

import (
	"github.com/s21platform/jarvis-bot/internal/config"
	"github.com/s21platform/jarvis-bot/internal/repository/postgres"
	"github.com/s21platform/jarvis-bot/internal/service/bot"
)

func main() {
	cfg := config.MustLoadConfig()
	db := postgres.New(cfg)
	Bot := bot.New(cfg, db)
	defer Bot.Close()
	Bot.Listen()
	select {}
}
