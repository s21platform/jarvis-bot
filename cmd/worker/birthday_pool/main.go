package main

import (
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/s21platform/jarvis-bot/internal/config"
	"github.com/s21platform/jarvis-bot/internal/repository/postgres"
	"github.com/s21platform/jarvis-bot/internal/repository/redis"
	worker "github.com/s21platform/jarvis-bot/internal/worker/birthday_pool"
)

func main() {
	cfg := config.MustLoadConfig()

	redis := redis.New(cfg)
	dbRepo := postgres.New(cfg)

	client := model.NewAPIv4Client(cfg.Bot.Url)
	client.SetOAuthToken(cfg.Bot.Token)

	worker := worker.NewWorker(redis, dbRepo, client, cfg)

	worker.Run()
}
