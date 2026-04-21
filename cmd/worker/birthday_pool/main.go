package main

import (
	"log"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/s21platform/jarvis-bot/internal/config"
	"github.com/s21platform/jarvis-bot/internal/repository/postgres"
	"github.com/s21platform/jarvis-bot/internal/repository/redis"
	worker "github.com/s21platform/jarvis-bot/internal/worker/birthday_pool"
)

// TODO: replace with real metrics implementation when metrics service is available
type noopMetrics struct{}

func (m *noopMetrics) Count(name string, value int64)        {}
func (m *noopMetrics) Disconnect()                           {}
func (m *noopMetrics) Duration(timestamp int64, name string) {}
func (m *noopMetrics) Gauge(name string, value float64)      {}
func (m *noopMetrics) Increment(name string)                 {}

func main() {
	cfg := config.MustLoadConfig()

	redis := redis.New(cfg)
	dbRepo := postgres.New(cfg)

	client := model.NewAPIv4Client(cfg.Bot.Url)
	client.SetOAuthToken(cfg.Bot.Token)

	w := worker.NewWorker(redis, dbRepo, client, cfg, &noopMetrics{})

	log.Println("Starting worker")
	w.Run()
	log.Println("worker shutting down")
}
