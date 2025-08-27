package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/s21platform/jarvis-bot/internal/config"
	"github.com/s21platform/jarvis-bot/internal/repository/postgres"
	"github.com/s21platform/jarvis-bot/internal/repository/redis"
	worker "github.com/s21platform/jarvis-bot/internal/worker/birthday_reminder"
)

type metrics struct{}

func (m *metrics) Increment(name string)                 {}
func (m *metrics) Count(name string, value int64)        {}
func (m *metrics) Duration(startTime int64, name string) {}

func main() {
	cfg := config.MustLoadConfig()

	redisRepo := redis.New(cfg)

	postgresRepo := postgres.New(cfg)

	client := model.NewAPIv4Client(cfg.Bot.Url)
	client.SetOAuthToken(cfg.Bot.Token)

	metrics := &metrics{}

	w := worker.NewWorker(redisRepo, postgresRepo, client, metrics)

	go w.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}
