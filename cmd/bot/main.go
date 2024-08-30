package main

import (
	"github.com/s21platform/jarvis-bot/internal/config"
	"github.com/s21platform/jarvis-bot/internal/service/bot"
)

func main() {
	cfg := config.MustLoadConfig()
	Bot := bot.New(cfg)
	defer Bot.Close()
	Bot.Listen()
	select {}
}
