package types

import (
	"github.com/mattermost/mattermost-server/v6/model"
)

// Command представляет собой интерфейс для всех команд бота
type Command interface {
	// Execute выполняет команду и возвращает сообщение для отправки
	Execute(args string) (string, error)
	// Name возвращает имя команды
	Name() string
	// Description возвращает описание команды
	Description() string
	// SetContext устанавливает контекст выполнения команды
	SetContext(ctx *CommandContext)
}

// BaseCommand содержит общие поля для всех команд
type BaseCommand struct {
	name        string
	description string
	context     *CommandContext
}

// NewBaseCommand создает новую базовую команду
func NewBaseCommand(name, description string) BaseCommand {
	return BaseCommand{
		name:        name,
		description: description,
	}
}

// CommandContext содержит контекст выполнения команды
type CommandContext struct {
	Post       *model.Post
	Channel    *model.Channel
	User       *model.User
	ProjectKey string // Ключ проекта Jira для текущего канала
}

// Name возвращает имя команды
func (c *BaseCommand) Name() string {
	return c.name
}

// Description возвращает описание команды
func (c *BaseCommand) Description() string {
	return c.description
}

// SetContext устанавливает контекст выполнения команды
func (c *BaseCommand) SetContext(ctx *CommandContext) {
	c.context = ctx
}

// GetContext возвращает контекст выполнения команды
func (c *BaseCommand) GetContext() *CommandContext {
	return c.context
}
