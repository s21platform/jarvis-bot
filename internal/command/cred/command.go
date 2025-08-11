package cred

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/s21platform/jarvis-bot/internal/pkg/types"
	"github.com/s21platform/jarvis-bot/internal/pkg/utils"
)

// Command реализует команду cred
type Command struct {
	types.BaseCommand
	db DbRepo
}

// NewCommand создает новую команду cred
func NewCommand(db DbRepo) *Command {
	return &Command{
		BaseCommand: types.NewBaseCommand("cred", "Показать учетные данные для сервиса"),
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

// Execute выполняет команду cred
func (c *Command) Execute(args string) (string, error) {
	if args == "" {
		return "Укажите название сервиса", nil
	}

	data, err := c.db.GetCred(context.Background(), args)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Sprintf("Учетные данные для сервиса '%s' не найдены", args), nil
		}
		return "", fmt.Errorf("не удалось получить учетные данные: %v", err)
	}

	rows := make([][]string, 0, len(data.Creds))
	for _, cred := range data.Creds {
		rows = append(rows, []string{cred.Name, cred.Value})
	}

	return utils.CreateTable([]string{"Имя", "Значение"}, rows), nil
}
