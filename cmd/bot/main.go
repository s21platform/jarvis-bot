package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/mattermost/mattermost-server/v6/model"

	"github.com/s21platform/jarvis-bot/internal/command"
	"github.com/s21platform/jarvis-bot/internal/config"
	"github.com/s21platform/jarvis-bot/internal/jira"
	"github.com/s21platform/jarvis-bot/internal/repository/postgres"
	"github.com/s21platform/jarvis-bot/internal/service/bot"

	"log"
)

// TODO: replace with real metrics implementation when metrics service is available
type noopMetrics struct{}

func (m *noopMetrics) Count(name string, value int64)        {}
func (m *noopMetrics) Disconnect()                           {}
func (m *noopMetrics) Duration(timestamp int64, name string) {}
func (m *noopMetrics) Gauge(name string, value float64)      {}
func (m *noopMetrics) Increment(name string)                 {}

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	cfg := config.MustLoadConfig()

	db := postgres.New(cfg)
	defer db.Close()

	jiraClient := jira.New(cfg)

	client := model.NewAPIv4Client(cfg.Bot.Url)
	client.SetOAuthToken(cfg.Bot.Token)

	user, _, err := client.GetMe("")
	if err != nil {
		log.Fatalf("Не удалось получить информацию о пользователе: %v", err)
	}

	cmdFactory := command.NewFactory(client, user, db, jiraClient)

	b := bot.New(cfg, cmdFactory, &noopMetrics{})
	defer b.Close()

	b.Listen()

	<-sigChan
	log.Println("Получен сигнал завершения, graceful shutdown...")
}
