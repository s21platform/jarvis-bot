package command

import (
	"context"

	"github.com/s21platform/jarvis-bot/internal/model"
)

type JiraClient interface {
	CreateIssue(title string, issueType string, projectKey string, labels []string) (string, error)
	GetBaseURL() string
}

type DbRepo interface {
	GetCred(ctx context.Context, key string) (model.Data, error)
	GetNews(ctx context.Context) ([]model.News, error)
	IsNewsExists(ctx context.Context, threadID string) (bool, error)
	AddNews(ctx context.Context, threadID, channel, messageID string) error
	DeleteNews(ctx context.Context, threadID string) error
	GetConfig(ctx context.Context, key string) (model.Config, error)
}
