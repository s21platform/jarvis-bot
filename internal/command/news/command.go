package news

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/s21platform/jarvis-bot/internal/pkg/types"
	"github.com/s21platform/jarvis-bot/internal/pkg/utils"
	"github.com/s21platform/jarvis-bot/internal/service/bot"
)

// Command реализует команду news
type Command struct {
	types.BaseCommand
	client *model.Client4
	db     bot.DbRepo
	botID  string // ID бота для реакций
}

const (
	newsReaction = "zap" // Emoji код для молнии ⚡
)

// NewCommand создает новую команду news
func NewCommand(client *model.Client4, db bot.DbRepo, botID string) *Command {
	return &Command{
		BaseCommand: types.NewBaseCommand("news", "Управление новостями (add/delete/show)"),
		client:      client,
		db:          db,
		botID:       botID,
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

// Execute выполняет команду news
func (c *Command) Execute(args string) (string, error) {
	ctx := context.Background()
	parts := strings.Fields(args)

	if len(parts) == 0 {
		return c.handleShow(ctx)
	}

	subcommand := parts[0]
	switch subcommand {
	case "add":
		if c.GetContext() == nil || c.GetContext().Post == nil {
			return "Эта команда должна быть выполнена в треде", nil
		}
		return c.handleAdd(ctx)
	case "delete":
		if c.GetContext() == nil || c.GetContext().Post == nil {
			return "Эта команда должна быть выполнена в треде", nil
		}
		return c.handleDelete(ctx)
	case "show":
		return c.handleShow(ctx)
	default:
		return "Неизвестная подкоманда. Используйте: add, delete или show", nil
	}
}

// handleAdd обрабатывает подкоманду add
func (c *Command) handleAdd(ctx context.Context) (string, error) {
	post := c.GetContext().Post
	channel := c.GetContext().Channel

	// Проверяем, что мы в треде
	if post.RootId == "" {
		return "Эта команда должна быть выполнена в треде", nil
	}

	exists, err := c.db.IsNewsExists(ctx, post.RootId)
	if err != nil {
		return "", fmt.Errorf("не удалось проверить существование новости: %v", err)
	}

	if exists {
		return "Эта новость уже добавлена", nil
	}

	// Получаем корневой пост
	rootPost, _, err := c.client.GetPost(post.RootId, "")
	if err != nil {
		return "", fmt.Errorf("не удалось получить корневой пост: %v", err)
	}

	err = c.db.AddNews(ctx, post.RootId, channel.DisplayName, rootPost.Id)
	if err != nil {
		return "", fmt.Errorf("не удалось добавить новость: %v", err)
	}

	// Добавляем реакцию к корневому посту
	reaction := &model.Reaction{
		UserId:    c.botID,
		PostId:    rootPost.Id,
		EmojiName: newsReaction,
	}
	_, _, err = c.client.SaveReaction(reaction)
	if err != nil {
		// Не возвращаем ошибку, так как добавление реакции не критично
		log.Printf("не удалось добавить реакцию: %v", err)
	}

	return "Новость успешно добавлена", nil
}

// handleDelete обрабатывает подкоманду delete
func (c *Command) handleDelete(ctx context.Context) (string, error) {
	post := c.GetContext().Post

	// Получаем информацию о новости перед удалением
	news, err := c.db.GetNews(ctx)
	if err != nil {
		return "", fmt.Errorf("не удалось получить информацию о новости: %v", err)
	}

	// Находим ID сообщения для удаления реакции
	var messageID string
	for _, n := range news {
		if n.ThreadID == post.RootId {
			messageID = n.MessageID
			break
		}
	}

	// Удаляем новость из БД
	err = c.db.DeleteNews(ctx, post.RootId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "Новость не найдена", nil
		}
		return "", fmt.Errorf("не удалось удалить новость: %v", err)
	}

	// Если нашли сообщение, удаляем реакцию
	if messageID != "" {
		reaction := &model.Reaction{
			UserId:    c.botID,
			PostId:    messageID,
			EmojiName: newsReaction,
		}
		_, err = c.client.DeleteReaction(reaction)
		if err != nil {
			// Не возвращаем ошибку, так как удаление реакции не критично
			log.Printf("не удалось удалить реакцию: %v", err)
		}
	}

	return "Новость успешно удалена", nil
}

// handleShow обрабатывает подкоманду show
func (c *Command) handleShow(ctx context.Context) (string, error) {
	news, err := c.db.GetNews(ctx)
	if err != nil {
		return "", fmt.Errorf("не удалось получить список новостей: %v", err)
	}

	if len(news) == 0 {
		return "Нет новостей за последние 3 недели", nil
	}

	rows := make([][]string, 0, len(news))
	for _, n := range news {
		// Получаем текст сообщения
		post, _, err := c.client.GetPost(n.MessageID, "")
		if err != nil {
			return "", fmt.Errorf("не удалось получить сообщение: %v", err)
		}

		// Форматируем сообщение: если длинное, то обрезаем и добавляем ссылку
		messageText := post.Message
		if len(messageText) > 80 {
			messageText = fmt.Sprintf("%s... [читать далее](https://mm.space-21.ru/space-21/pl/%s)", messageText[:77], n.MessageID)
		}

		// Форматируем дату в читаемый вид
		date := n.CreatedAt.Format("02.01.2006 15:04")
		rows = append(rows, []string{date, n.Channel, messageText})
	}

	return utils.CreateTable([]string{"Дата", "Канал", "Сообщение"}, rows), nil
}
