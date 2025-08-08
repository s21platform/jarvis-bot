package main

import (
	"log"

	"github.com/mattermost/mattermost-server/v6/model"

	"github.com/s21platform/jarvis-bot/internal/command"
	"github.com/s21platform/jarvis-bot/internal/config"
	"github.com/s21platform/jarvis-bot/internal/jira"
	"github.com/s21platform/jarvis-bot/internal/repository/postgres"
	"github.com/s21platform/jarvis-bot/internal/service/bot"
)

func main() {
	cfg := config.MustLoadConfig()

	// Инициализируем подключение к базе данных
	db := postgres.New(cfg)
	defer db.Close()

	jiraClient := jira.New(cfg)
	jiraClient.GetProjects()

	// Создаем клиента Mattermost и получаем информацию о пользователе
	client := model.NewAPIv4Client(cfg.Url)
	client.SetOAuthToken(cfg.Token)

	user, _, err := client.GetMe("")
	if err != nil {
		log.Fatalf("Не удалось получить информацию о пользователе: %v", err)
	}

	// Создаем фабрику команд
	cmdFactory := command.NewFactory(client, user, db)

	// Создаем и запускаем бота
	b := bot.New(cfg, cmdFactory)
	defer b.Close()
	b.Listen()
	select {}
}
