package main

import (
	"log"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/s21platform/jarvis-bot/internal/config"
	"github.com/s21platform/jarvis-bot/internal/repository/postgres"
	"github.com/s21platform/jarvis-bot/internal/repository/redis"
	worker "github.com/s21platform/jarvis-bot/internal/worker/birthday_pool"
	"github.com/s21platform/metrics-lib/pkg"
)

func main() {
	cfg := config.MustLoadConfig()

	redis := redis.New(cfg)
	dbRepo := postgres.New(cfg)

	client := model.NewAPIv4Client(cfg.Bot.Url)
	client.SetOAuthToken(cfg.Bot.Token)

	metrics, err := pkg.NewMetrics(cfg.Metrics.Host, cfg.Metrics.Port, "jarvis-bot", cfg.Platform.Env)
	if err != nil {
		log.Fatalf("Failed to create metrics: %v", err)
	}

	worker := worker.NewWorker(redis, dbRepo, client, cfg, metrics)

	log.Println("Starting worker")
	worker.Run()
	log.Println("worker shutting down")
}
