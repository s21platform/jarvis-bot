package bot

import (
	"context"

	"github.com/s21platform/jarvis-bot/internal/model"
	"github.com/s21platform/jarvis-bot/internal/pkg/types"
)

// DbRepo представляет собой интерфейс для работы с базой данных
type DbRepo interface {
	// GetCred возвращает учетные данные для сервиса
	GetCred(ctx context.Context, name string) (model.Data, error)

	// AddNews добавляет новость в базу данных
	AddNews(ctx context.Context, threadID, channel, messageID string) error

	// DeleteNews удаляет новость из базы данных
	DeleteNews(ctx context.Context, threadID string) error

	// GetNews возвращает список всех новостей за последние 3 недели
	GetNews(ctx context.Context) ([]model.News, error)

	// IsNewsExists проверяет существование новости в базе данных
	IsNewsExists(ctx context.Context, threadID string) (bool, error)

	// Close закрывает соединение с базой данных
	Close()
}

// CommandFactory представляет собой интерфейс для работы с командами
type CommandFactory interface {
	// GetCommand возвращает команду по имени
	GetCommand(name string) types.Command

	// GetAllCommands возвращает список всех команд
	GetAllCommands() []types.Command
}
