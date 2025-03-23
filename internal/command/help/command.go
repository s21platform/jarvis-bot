package help

import (
	"fmt"
	"strings"

	"github.com/s21platform/jarvis-bot/internal/pkg/types"
)

// Command реализует команду help
type Command struct {
	types.BaseCommand
	commands []types.Command
}

// NewCommand создает новую команду help
func NewCommand(commands []types.Command) *Command {
	return &Command{
		BaseCommand: types.NewBaseCommand("help", "Показать список доступных команд"),
		commands:    commands,
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

// Execute выполняет команду help
func (c *Command) Execute(args string) (string, error) {
	var sb strings.Builder
	sb.WriteString("Доступные команды:\n\n")

	for _, cmd := range c.commands {
		sb.WriteString(fmt.Sprintf("**%s** - %s\n", cmd.Name(), cmd.Description()))
	}

	return sb.String(), nil
}
