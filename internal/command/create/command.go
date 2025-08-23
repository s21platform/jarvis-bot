package create

import (
	"fmt"
	"strings"

	"github.com/s21platform/jarvis-bot/internal/pkg/types"
)

// Command реализует команду create
type Command struct {
	types.BaseCommand
	db DbRepo
	jC JiraClient
}

// GetContext возвращает контекст команды
func (c *Command) GetContext() *types.CommandContext {
	return c.BaseCommand.GetContext()
}

// NewCommand создает новую команду create
func NewCommand(db DbRepo, jC JiraClient) *Command {
	return &Command{
		BaseCommand: types.NewBaseCommand("create", "Создать задачу или баг (task/bug)"),
		db:          db,
		jC:          jC,
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
	ctx := c.GetContext()
	if ctx == nil || ctx.Channel == nil {
		return "", fmt.Errorf("не удалось получить контекст команды")
	}

	if ctx.ProjectKey == "" {
		return fmt.Sprintf("Канал %s не привязан к проекту Jira", ctx.Channel.Name), nil
	}

	var labels []string
	if ctx.Label != nil {
		labels = []string{*ctx.Label}
	}

	key, err := c.jC.CreateIssue(title, "Task", ctx.ProjectKey, labels)
	if err != nil {
		return "", fmt.Errorf("failed to create jira task: %w", err)
	}

	return fmt.Sprintf("Создана задача [%s](%s/browse/%s) :bulb:", key, c.jC.GetBaseURL(), key), nil
}

// handleBug обрабатывает создание бага
func (c *Command) handleBug(title string) (string, error) {
	ctx := c.GetContext()
	if ctx == nil || ctx.Channel == nil {
		return "", fmt.Errorf("не удалось получить контекст команды")
	}

	if ctx.ProjectKey == "" {
		return fmt.Sprintf("Канал %s не привязан к проекту Jira", ctx.Channel.Name), nil
	}

	var labels []string
	if ctx.Label != nil {
		labels = []string{*ctx.Label}
	}

	key, err := c.jC.CreateIssue(title, "Bug", ctx.ProjectKey, labels)
	if err != nil {
		return "", fmt.Errorf("failed to create jira bug: %w", err)
	}

	return fmt.Sprintf("Создан баг [%s](%s/browse/%s) :ladybug:", key, c.jC.GetBaseURL(), key), nil
}
