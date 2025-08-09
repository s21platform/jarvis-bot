package create

import (
	"fmt"
	"strings"

	"github.com/s21platform/jarvis-bot/internal/pkg/types"
	"github.com/s21platform/jarvis-bot/internal/service/bot"
)

// Command реализует команду create
type Command struct {
	types.BaseCommand
	db bot.DbRepo
}

// NewCommand создает новую команду create
func NewCommand(db bot.DbRepo) *Command {
	return &Command{
		BaseCommand: types.NewBaseCommand("create", "Создать задачу или баг (task/bug)"),
		db:          db,
	}
}

// Name возвращает имя команды
func (c *Command) Name() string {
	return c.BaseCommand.Name()
}

// Description возвращает описание команды
func (c *Command) Description() string {
	return c.BaseCommand.Description()
}

// Execute выполняет команду create
func (c *Command) Execute(args string) (string, error) {
	parts := strings.Fields(args)
	if len(parts) < 2 {
		return "Использование: create task|bug <заголовок>", nil
	}

	subcommand := parts[0]
	title := strings.Join(parts[1:], " ")

	switch subcommand {
	case "task":
		return c.handleTask(title)
	case "bug":
		return c.handleBug(title)
	default:
		return "Неизвестная подкоманда. Используйте: task или bug", nil
	}
}

// handleTask обрабатывает создание задачи
func (c *Command) handleTask(title string) (string, error) {
	// TODO: Реализовать создание задачи
	return fmt.Sprintf("Создание задачи с заголовком: %s", title), nil
}

// handleBug обрабатывает создание бага
func (c *Command) handleBug(title string) (string, error) {
	// TODO: Реализовать создание бага
	return fmt.Sprintf("Создание бага с заголовком: %s", title), nil
}
