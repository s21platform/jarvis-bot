package command

import (
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/s21platform/jarvis-bot/internal/command/create"
	"github.com/s21platform/jarvis-bot/internal/command/cred"
	"github.com/s21platform/jarvis-bot/internal/command/help"
	"github.com/s21platform/jarvis-bot/internal/command/news"
	"github.com/s21platform/jarvis-bot/internal/pkg/types"
	"github.com/s21platform/jarvis-bot/internal/service/bot"
)

// Factory создает и управляет командами
type Factory struct {
	commands map[string]types.Command
}

// NewFactory создает новую фабрику команд
func NewFactory(client *model.Client4, user *model.User, db bot.DbRepo, jC JiraClient) *Factory {
	f := &Factory{
		commands: make(map[string]types.Command),
	}

	// Создаем команды
	credCmd := cred.NewCommand(db)
	newsCmd := news.NewCommand(client, db, user.Id)
	createCmd := create.NewCommand(db, jC)

	// Добавляем команды в map
	commands := []types.Command{
		credCmd,
		newsCmd,
		createCmd,
	}

	// Создаем help команду со списком всех команд
	helpCmd := help.NewCommand(commands)
	commands = append(commands, helpCmd)

	// Заполняем map команд
	for _, cmd := range commands {
		f.commands[cmd.Name()] = cmd
	}

	return f
}

// GetCommand возвращает команду по имени
func (f *Factory) GetCommand(name string) types.Command {
	return f.commands[name]
}

// GetAllCommands возвращает список всех команд
func (f *Factory) GetAllCommands() []types.Command {
	commands := make([]types.Command, 0, len(f.commands))
	for _, cmd := range f.commands {
		commands = append(commands, cmd)
	}
	return commands
}
