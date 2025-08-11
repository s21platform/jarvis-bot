package news

import (
	"context"

	"github.com/s21platform/jarvis-bot/internal/model"
)

type DbRepo interface {
	GetNews(ctx context.Context) ([]model.News, error)
	IsNewsExists(ctx context.Context, threadID string) (bool, error)
	AddNews(ctx context.Context, threadID, channel, messageID string) error
	DeleteNews(ctx context.Context, threadID string) error
}